#!/bin/bash

# Update Protos Script
# This script handles updating proto submodules, generating Go code, and managing versions
# Can be invoked by GitHub Actions or run locally

set -euo pipefail

# Default values
RELEASE_VERSION=""
SDK_VERSION=""
BRANCH_NAME=""
PR_TITLE=""
DRY_RUN=false
VERBOSE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Usage function
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Update proto submodules, generate Go code, and manage versions.

OPTIONS:
    -r, --release-version VERSION    Cloud API release version to checkout (e.g., v1.2.3)
    -s, --sdk-version VERSION        SDK version to set (e.g., 0.7.0) (optional)
    -b, --branch-name NAME          Branch name for the PR (default: auto-generated)
    -t, --pr-title TITLE            PR title (default: auto-generated)
    -d, --dry-run                   Show what would be done without making changes
    -v, --verbose                   Enable verbose output
    -h, --help                      Show this help message

EXAMPLES:
    # Update to latest proto version without changing SDK version
    $0

    # Update to latest proto version with specific SDK version
    $0 --sdk-version 0.7.0

    # Update to specific release version with SDK version
    $0 --release-version v1.2.3 --sdk-version 0.7.0

    # Dry run to see what would happen
    $0 --sdk-version 0.7.0 --dry-run --verbose

    # Update with custom branch and PR title
    $0 --sdk-version 0.7.0 --branch-name my-update --pr-title "Custom Update"

ENVIRONMENT VARIABLES:
    GITHUB_TOKEN                    GitHub token for creating PRs (required for PR creation)
    GITHUB_REPOSITORY              Repository in format owner/repo (required for PR creation)
    GITHUB_ACTIONS                 Set to 'true' when running in GitHub Actions

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -r|--release-version)
                RELEASE_VERSION="$2"
                # Allow empty value to fallback to no set
                if [[ -n "$RELEASE_VERSION" ]]; then
                    if ! validate_version "$RELEASE_VERSION"; then
                        log_error "Invalid release version format: $RELEASE_VERSION"
                        log_error "Expected format: v1.2.3 or 1.2.3 (with optional pre-release suffix)"
                        exit 1
                    fi
                fi
                shift 2
                ;;
            -s|--sdk-version)
                SDK_VERSION="$2"
                # Allow empty value to fallback to no set
                if [[ -n "$SDK_VERSION" ]]; then
                    if ! validate_version "$SDK_VERSION"; then
                        log_error "Invalid SDK version format: $SDK_VERSION"
                        log_error "Expected format: v1.2.3 or 1.2.3 (with optional pre-release suffix)"
                        exit 1
                    fi
                fi
                shift 2
                ;;
            -b|--branch-name)
                BRANCH_NAME="$2"
                # Allow empty value to fallback to no set (auto-generated)
                if [[ -n "$BRANCH_NAME" ]]; then
                    if ! validate_branch_name "$BRANCH_NAME"; then
                        log_error "Invalid branch name format: $BRANCH_NAME"
                        log_error "Branch name must be alphanumeric with dashes, underscores, slashes, or dots"
                        exit 1
                    fi
                fi
                shift 2
                ;;
            -t|--pr-title)
                PR_TITLE="$2"
                # Allow empty value to fallback to no set (auto-generated)
                # Validate PR title doesn't contain dangerous characters if provided
                if [[ -n "$PR_TITLE" ]] && [[ "$PR_TITLE" =~ [\`\$\(\)] ]]; then
                    log_error "Invalid PR title: contains dangerous characters"
                    exit 1
                fi
                shift 2
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

# Check if running in GitHub Actions
is_github_actions() {
    [[ "${GITHUB_ACTIONS:-}" == "true" ]]
}

# Check if we can create PRs
can_create_pr() {
    [[ -n "${GITHUB_TOKEN:-}" && -n "${GITHUB_REPOSITORY:-}" ]]
}

# Validate version string format (semantic versioning with optional 'v' prefix)
# Returns 0 if valid, 1 if invalid
validate_version() {
    local version="$1"
    # Allow: v1.2.3, 1.2.3, 1.2.3-alpha, 1.2.3-rc.1, etc.
    if [[ "$version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9._-]+)?$ ]]; then
        return 0
    else
        return 1
    fi
}

# Validate branch name format
# Returns 0 if valid, 1 if invalid
validate_branch_name() {
    local branch="$1"
    # Git branch names: alphanumeric, dash, underscore, slash, dot (no spaces, special chars)
    # Cannot start with dot, dash, or end with .lock
    if [[ "$branch" =~ ^[a-zA-Z0-9][a-zA-Z0-9/_.-]*$ ]] && \
       [[ ! "$branch" =~ \.\.|\.\.|@\{|\\|\ |\~|\^|\:|\?|\*|\[|//|\.lock$ ]]; then
        return 0
    else
        return 1
    fi
}

# Sanitize string for use in sed patterns (escape special sed characters)
sanitize_for_sed() {
    local string="$1"
    # Escape forward slashes, ampersands, backslashes, and newlines
    printf '%s' "$string" | sed 's/[\/&\\]/\\&/g'
}

# Cross-platform sed in-place editing
# Usage: sed_inplace 's/pattern/replacement/' file
sed_inplace() {
    local pattern="$1"
    local file="$2"
    
    # macOS (BSD sed) requires an extension argument, Linux (GNU sed) does not
    if [[ "$(uname)" == "Darwin" ]]; then
        sed -i '' "$pattern" "$file"
    else
        sed -i "$pattern" "$file"
    fi
}

# Update proto submodule
update_proto_submodule() {
    log_info "Updating proto submodule..."
    
    if [[ -n "$RELEASE_VERSION" ]]; then
        log_info "Checking out specific release version: $RELEASE_VERSION"
        
        # Update submodule to latest first
        git submodule update --init --recursive
        
        # Navigate to submodule and checkout specific release
        cd proto/cloud-api
        git fetch --tags
        git checkout "$RELEASE_VERSION"
        cd ../..
        
        # Get the version from the checked out release
        PROTO_VERSION=$(cat proto/cloud-api/VERSION)
        if ! validate_version "$PROTO_VERSION"; then
            log_error "Invalid proto version format in VERSION file: $PROTO_VERSION"
            exit 1
        fi
        log_info "Checked out proto version: $PROTO_VERSION"
    else
        log_info "Updating proto submodule to get latest changes..."
        git submodule update --recursive --remote --merge
        
        # Get the current proto version
        PROTO_VERSION=$(cat proto/cloud-api/VERSION)
        if ! validate_version "$PROTO_VERSION"; then
            log_error "Invalid proto version format in VERSION file: $PROTO_VERSION"
            exit 1
        fi
        log_info "Current proto version: $PROTO_VERSION"
    fi
    
    # Check if there are any changes
    if git diff --quiet proto/cloud-api; then
        log_warning "No changes in proto submodule"
        return 1
    else
        log_success "Changes detected in proto submodule"
        return 0
    fi
}

# Generate Go code from protos
generate_go_code() {
    log_info "Generating Go code from protos..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would install dependencies..."
        log_info "[DRY RUN] Would clean generated files..."
        log_info "[DRY RUN] Would generate Go code from protos..."
        log_info "[DRY RUN] Would move generated files to correct location..."
        log_info "[DRY RUN] Would update default API version in cloudclient/options.go to $PROTO_VERSION"
        if [[ -n "$SDK_VERSION" ]]; then
            log_info "[DRY RUN] Would update SDK version in cloudclient/options.go to $SDK_VERSION"
            # Export the provided SDK version for dry run
            export SDK_VERSION="$SDK_VERSION"
        fi
        return 0
    fi
    
    log_info "Installing dependencies..."
    go generate ./internal/build/tools.go
    
    log_info "Cleaning generated files..."
    rm -rf api/*
    
    log_info "Generating Go code from protos..."
    buf generate
    
    log_info "Moving generated files to correct location..."
    mv -f api/temporal/api/cloud/* api && rm -rf api/temporal
    
    log_info "Updating default API version in cloudclient/options.go to $PROTO_VERSION..."
    # Sanitize version strings for use in sed
    PROTO_VERSION_SAFE=$(sanitize_for_sed "$PROTO_VERSION")
    sed_inplace "s/defaultAPIVersion = \".*\"/defaultAPIVersion = \"$PROTO_VERSION_SAFE\"/" cloudclient/options.go
    
    # Only update SDK version if provided
    if [[ -n "$SDK_VERSION" ]]; then
        log_info "Updating SDK version in cloudclient/options.go to $SDK_VERSION..."
        # Extract current SDK version for logging
        CURRENT_SDK_VERSION=$(grep 'sdkVersion.*=' cloudclient/options.go | sed 's/.*= \"\(.*\)\".*/\1/')
        log_info "Current SDK version: $CURRENT_SDK_VERSION"
        log_info "New SDK version: $SDK_VERSION"
        
        # Update SDK version in options.go
        SDK_VERSION_SAFE=$(sanitize_for_sed "$SDK_VERSION")
        sed_inplace "s/sdkVersion.*= \".*\"/sdkVersion        = \"$SDK_VERSION_SAFE\"/" cloudclient/options.go
        
        # Export variables for later use
        export SDK_VERSION="$SDK_VERSION"
    fi
}

# Check for generated changes
check_generated_changes() {
    log_info "Checking for generated changes..."
    
    if git diff --quiet; then
        log_warning "No changes to commit after proto generation"
        return 1
    else
        log_success "Generated changes detected"
        return 0
    fi
}

# Create branch and commit changes
create_branch_and_commit() {
    log_info "Creating branch and committing changes..."
    
    # Generate branch name if not provided
    if [[ -z "$BRANCH_NAME" ]]; then
        # Create a clean version string (remove 'v' prefix and replace dots with dashes)
        PROTO_VERSION_CLEAN="${PROTO_VERSION#v}"
        PROTO_VERSION_CLEAN="${PROTO_VERSION_CLEAN//./-}"
        
        # Generate a random suffix to ensure uniqueness
        RANDOM_SUFFIX=$(date +%s%5n | cut -b2-10)
        
        if [[ -n "$SDK_VERSION" ]]; then
            # Include both proto and SDK versions in branch name
            SDK_VERSION_CLEAN="${SDK_VERSION#v}"
            SDK_VERSION_CLEAN="${SDK_VERSION_CLEAN//./-}"
            BRANCH_NAME="update-protos-${PROTO_VERSION_CLEAN}-sdk-${SDK_VERSION_CLEAN}-${RANDOM_SUFFIX}"
        else
            # Include only proto version in branch name
            BRANCH_NAME="update-protos-${PROTO_VERSION_CLEAN}-${RANDOM_SUFFIX}"
        fi
        
        # Validate the generated branch name
        if ! validate_branch_name "$BRANCH_NAME"; then
            log_error "Generated branch name is invalid: $BRANCH_NAME"
            exit 1
        fi
    fi
    log_info "Branch name: $BRANCH_NAME"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would configure git"
        log_info "[DRY RUN] Would create and checkout new branch: $BRANCH_NAME"
        log_info "[DRY RUN] Would add all changes"
        log_info "[DRY RUN] Would commit changes"
        log_info "[DRY RUN] Would push branch"
        
        # Set dummy commit message for dry run
        if [[ -n "$SDK_VERSION" ]]; then
            export COMMIT_MSG="[DRY RUN] Update protos to version $PROTO_VERSION and set SDK to $SDK_VERSION"
        else
            export COMMIT_MSG="[DRY RUN] Update protos to version $PROTO_VERSION"
        fi
        return 0
    fi
    
    # Configure git
    git config --local user.email "action@github.com"
    git config --local user.name "GitHub Action"
    
    # Create and checkout new branch
    git checkout -b "$BRANCH_NAME"
    
    # Add all changes
    git add .
    
    # Generate commit message
    if [[ -n "$SDK_VERSION" ]]; then
        if [[ -n "$RELEASE_VERSION" ]]; then
            COMMIT_MSG="Update protos to release $RELEASE_VERSION (version $PROTO_VERSION) and set SDK to $SDK_VERSION"
        else
            COMMIT_MSG="Update protos to version $PROTO_VERSION and set SDK to $SDK_VERSION"
        fi
    else
        if [[ -n "$RELEASE_VERSION" ]]; then
            COMMIT_MSG="Update protos to release $RELEASE_VERSION (version $PROTO_VERSION)"
        else
            COMMIT_MSG="Update protos to version $PROTO_VERSION"
        fi
    fi
    log_info "Commit message: $COMMIT_MSG"
    
    # Commit changes
    git commit -m "$COMMIT_MSG"
    
    # Push branch
    git push origin "$BRANCH_NAME"
    
    # Export variables for later use
    export COMMIT_MSG
}

# Create pull request
create_pull_request() {
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY RUN] Would create pull request"
        return 0
    fi
    
    if ! can_create_pr; then
        log_warning "Cannot create PR: Missing GITHUB_TOKEN or GITHUB_REPOSITORY"
        return 0
    fi
    
    log_info "Creating pull request..."
    
    # Generate PR title if not provided
    if [[ -z "$PR_TITLE" ]]; then
        if [[ -n "$SDK_VERSION" ]]; then
            if [[ -n "$RELEASE_VERSION" ]]; then
                PR_TITLE="Update protos to release $RELEASE_VERSION (version $PROTO_VERSION) and set SDK to $SDK_VERSION"
            else
                PR_TITLE="Update protos to version $PROTO_VERSION and set SDK to $SDK_VERSION"
            fi
        else
            if [[ -n "$RELEASE_VERSION" ]]; then
                PR_TITLE="Update protos to release $RELEASE_VERSION (version $PROTO_VERSION)"
            else
                PR_TITLE="Update protos to version $PROTO_VERSION"
            fi
        fi
    fi
    
    # Generate PR body
    if [[ -n "$SDK_VERSION" ]]; then
        PR_BODY="## Changes

This PR updates the proto files to version **$PROTO_VERSION** and sets the SDK version to **$SDK_VERSION**."
    else
        PR_BODY="## Changes

This PR updates the proto files to version **$PROTO_VERSION**."
    fi

    if [[ -n "$RELEASE_VERSION" ]]; then
        PR_BODY+="

**Release:** $RELEASE_VERSION"
    fi
    
    PR_BODY+="

### What was done:
- Updated proto submodule${RELEASE_VERSION:+ to release $RELEASE_VERSION}
- Regenerated Go code from proto files
- Updated default API version in \`cloudclient/options.go\` to \`$PROTO_VERSION\`"

    if [[ -n "$SDK_VERSION" ]]; then
        PR_BODY+="
- Set SDK version in \`cloudclient/options.go\` to \`$SDK_VERSION\`"
    fi

    PR_BODY+="

### Generated files:
- All files in \`api/\` directory have been regenerated
- API version updated to \`$PROTO_VERSION\`"

    if [[ -n "$SDK_VERSION" ]]; then
        PR_BODY+="
- SDK version updated to \`$SDK_VERSION\`"
    fi

    PR_BODY+="

This PR was automatically created by the Update Protos workflow."

    # Check if PR already exists
    EXISTING_PR=$(gh pr list --head "$BRANCH_NAME" --state open --json number --jq '.[0].number' 2>/dev/null || echo "")
    
    if [[ -n "$EXISTING_PR" ]]; then
        log_warning "PR already exists: #$EXISTING_PR"
        return 0
    fi
    
    # Create PR
    PR_URL=$(gh pr create --title "$PR_TITLE" --body "$PR_BODY" --head "$BRANCH_NAME" --base main 2>/dev/null || echo "")
    
    if [[ -n "$PR_URL" ]]; then
        log_success "Created PR: $PR_URL"
    else
        log_error "Failed to create PR"
        return 1
    fi
}

# Print summary
print_summary() {
    echo
    log_info "=== SUMMARY ==="
    
    if [[ "$NO_CHANGES" == "true" ]]; then
        log_success "No changes detected in proto submodule. No action needed."
    elif [[ "$NO_GENERATED_CHANGES" == "true" ]]; then
        log_success "Proto submodule updated but no generated code changes detected."
    else
        log_success "Successfully updated protos:"
        echo "  - Proto version: $PROTO_VERSION"
        if [[ -n "$SDK_VERSION" ]]; then
            echo "  - SDK version: $SDK_VERSION"
        fi
        if [[ -n "$RELEASE_VERSION" ]]; then
            echo "  - Release: $RELEASE_VERSION"
        fi
        echo "  - Branch: $BRANCH_NAME"
        echo "  - Commit: $COMMIT_MSG"
    fi
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    # Set verbose mode
    if [[ "$VERBOSE" == "true" ]]; then
        set -x
    fi
    
    # Initialize variables
    NO_CHANGES=false
    NO_GENERATED_CHANGES=false
    
    log_info "Starting proto update process..."
    
    # Update proto submodule
    if ! update_proto_submodule; then
        NO_CHANGES=true
        print_summary
        exit 0
    fi
    
    # Generate Go code
    generate_go_code
    
    # Check for generated changes
    if ! check_generated_changes; then
        NO_GENERATED_CHANGES=true
        print_summary
        exit 0
    fi
    
    # Create branch and commit changes
    create_branch_and_commit
    
    # Create pull request (only if not dry run and we can create PRs)
    if [[ "$DRY_RUN" != "true" && "$NO_CHANGES" != "true" && "$NO_GENERATED_CHANGES" != "true" ]]; then
        create_pull_request
    fi
    
    # Print summary
    print_summary
}

# Run main function with all arguments
main "$@"

#!/bin/bash

# Version update script for Flow-Sight application
# Usage: ./scripts/update-version.sh <new_version>
# Example: ./scripts/update-version.sh 1.1.0

set -e

# Check if version argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <new_version>"
    echo "Example: $0 1.1.0"
    exit 1
fi

NEW_VERSION="$1"
OLD_VERSION=$(grep -o '"version": "[^"]*"' frontend/package.json | cut -d'"' -f4)

echo "Updating version from $OLD_VERSION to $NEW_VERSION"

# Function to update version in a file
update_version() {
    local file="$1"
    local pattern="$2"
    local replacement="$3"
    
    if [ -f "$file" ]; then
        echo "Updating $file..."
        sed -i.bak "$pattern" "$file"
        rm "$file.bak"
        echo "✓ Updated $file"
    else
        echo "⚠️  File not found: $file"
    fi
}

# Update backend version
echo "📦 Updating backend version..."
update_version "backend/internal/version/version.go" "s/Version = \"[^\"]*\"/Version = \"$NEW_VERSION\"/g"

# Update frontend package.json
echo "🌐 Updating frontend version..."
update_version "frontend/package.json" "s/\"version\": \"[^\"]*\"/\"version\": \"$NEW_VERSION\"/g"

# Update helm chart appVersion only
echo "⚓ Updating helm chart appVersion..."
update_version "helm-chart/Chart.yaml" "s/appVersion: \"[^\"]*\"/appVersion: \"$NEW_VERSION\"/g"

# Update helm chart values files - only backend and frontend image tags
echo "🏷️  Updating helm chart image tags..."

# Function to update specific image tags in YAML files
update_image_tag() {
    local file="$1"
    local service="$2"
    local new_version="$3"
    
    if [ -f "$file" ]; then
        echo "Updating $service image tag in $file..."
        # Use awk to update only the tag under the specific service section
        awk -v service="$service" -v new_version="$new_version" '
        BEGIN { in_section = 0 }
        /^[a-zA-Z]/ { in_section = 0 }
        $0 ~ "^" service ":" { in_section = 1 }
        in_section && /^  image:/ { in_image = 1; print; next }
        in_section && in_image && /^    tag:/ { 
            gsub(/"[^"]*"/, "\"" new_version "\""); 
            in_image = 0 
        }
        { print }
        ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
        echo "✓ Updated $service tag in $file"
    else
        echo "⚠️  File not found: $file"
    fi
}

# Update backend and frontend image tags
update_image_tag "helm-chart/values.yaml" "backend" "$NEW_VERSION"
update_image_tag "helm-chart/values.yaml" "frontend" "$NEW_VERSION"
update_image_tag "helm-chart/values-pke.yaml" "backend" "$NEW_VERSION"
update_image_tag "helm-chart/values-pke.yaml" "frontend" "$NEW_VERSION"

echo ""
echo "✅ Version update completed!"
echo "📋 Summary of changes:"
echo "   Backend: backend/internal/version/version.go"
echo "   Frontend: frontend/package.json"
echo "   Helm Chart: helm-chart/Chart.yaml"
echo "   Helm Values: helm-chart/values.yaml, helm-chart/values-pke.yaml"
echo ""
echo "🔍 You can verify the changes with:"
echo "   git diff"
echo ""
echo "💡 Don't forget to:"
echo "   1. Test the application"
echo "   2. Update CHANGELOG.md if you have one"
echo "   3. Create a git tag: git tag v$NEW_VERSION"
echo "   4. Commit and push changes"

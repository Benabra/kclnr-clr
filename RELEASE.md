# Creating a New Release for kclnr-cli

This document provides step-by-step instructions for creating a new release of the `kclnr-cli` project.

## Steps to Create a New Release

1. **Ensure All Changes are Committed**

   Make sure all your changes are committed before you start the release process. You can check the status of your changes with:

   ```sh
   git status
   ```

2. **Run the Release Command**

   To create a new release, run the release command. This command will perform the following actions:
   - Tidy the Go modules
   - Build the project
   - Get the version from the built binary
   - Prompt for a commit message
   - Add all changes
   - Commit the changes
   - Push the changes to the `main` branch
   - Create a new tag
   - Push the new tag to the remote repository

   ```sh
   go run main.go release
   ```

3. **Enter the Commit Message**

   When prompted, enter a commit message that describes the changes included in this release.

   ```
   Enter commit message: Describe the changes in this release
   ```

4. **Verify the Release**

   After running the release command, verify that the new tag has been created and pushed to the remote repository. You can check the tags with:

   ```sh
   git tag
   ```

   And verify that the tag exists on the remote repository:

   ```sh
   git ls-remote --tags origin
   ```

## Example

Here is an example of creating a new release:

1. Ensure all changes are committed:

   ```sh
   git status
   ```

2. Run the release command:

   ```sh
   go run main.go release
   ```

3. Enter the commit message when prompted:

   ```
   Enter commit message: Fix module import paths and update project structure
   ```

4. Verify the new tag:

   ```sh
   git tag
   git ls-remote --tags origin
   ```

## Notes

- Make sure you have the necessary permissions to push changes and tags to the remote repository.
- Ensure that the `go.mod` file and other necessary files are correctly set up before running the release command.

By following these steps, you can create a new release for the `kclnr-cli` project efficiently and ensure that the release process is consistent and well-documented.
```

### Adding `RELEASE.md` to Your Repository

1. Create the `RELEASE.md` file in the root directory of your project:

   ```sh
   touch RELEASE.md
   ```

2. Copy and paste the content provided above into the `RELEASE.md` file.

3. Add and commit the `RELEASE.md` file to your repository:

   ```sh
   git add RELEASE.md
   git commit -m "Add RELEASE.md to explain how to create a new release"
   git push origin main
   ```

By adding this `RELEASE.md` file to your project, you provide clear and detailed instructions on how to create a new release, ensuring that anyone contributing to the project can follow the same process.
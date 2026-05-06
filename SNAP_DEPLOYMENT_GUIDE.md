# Snap Deployment Guide (Clonar Mídia)

This document explains how to maintain the automated deployment of Clonar Mídia to the Snap Store.

## 1. Authentication (GitHub Secrets)

The GitHub Action uses a login token to publish the Snap package. This token is stored in the repository secrets.

**If the build fails with `login_data is empty`:**

1.  **Export a new login token** on your local machine:
    ```bash
    snapcraft export-login --snaps=clonarmidia --expires=2026-12-31 snap.login
    ```
2.  **Copy the content** of the generated `snap.login` file.
3.  **Update GitHub Secret**:
    - Go to your repository on GitHub.
    - Navigate to **Settings** > **Secrets and variables** > **Actions**.
    - Update or create a secret named: `SNAPCRAFT_STORE_CREDENTIALS`.
    - Paste the content of the token.

## 2. The "Demo" Strategy

To ensure **instant approval** in the Snap Store (avoiding manual human review that can take days), we use a "Demo Strategy":

-   **Confinement**: We use `strict` confinement.
-   **Permissions**: We **do not** request privileged plugs like `block-devices` or `udisks2` in the Snap version.
-   **UX**: The application detects it's running in a restricted Snap sandbox and displays a banner inviting the user to download the **Full Version** (.deb/.rpm) from GitHub.

## 3. Manual Release (Emergency)

If GitHub Actions fails and you need to push a version manually:

1.  Build the snap locally:
    ```bash
    snapcraft
    ```
2.  Push to the store:
    ```bash
    snapcraft upload --release=edge clonarmidia_1.0.0_amd64.snap
    ```

## 4. Full Version Distribution

The full version (with total hardware access) is distributed via the `build/bin` folder or GitHub Releases. 
- Ensure the binaries are updated in `build/bin` before pushing to GitHub.
- The Snap version's "Download Full Version" button points directly to: `https://github.com/erascardsilva/clonarMidia/tree/main/build/bin`

---
**Erasmo Cardoso**  
Software Engineer | Electronics Technician

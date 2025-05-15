# Folder Management GUI

This is a web-based folder management GUI built with Go and Tailwind CSS. It allows you to browse, create, upload, and delete folders and files on your local system via a browser interface. This README provides instructions on running the application, changing the default path, and exposing it to other devices using ngrok.


---

## 📸 Preview

![Remote File Manager](screenshot-remore-file-access.jpg)

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.16 or later): [Download and install Go](https://go.dev/doc/install)
- **Git**: To clone the repository (optional if you already have the code)
- A modern web browser (e.g., Chrome, Firefox)
- An internet connection (for downloading ngrok and accessing the app from other devices)

## Running the Application

Follow these steps to run the application locally on your machine:

1. **Clone the Repository** (if you haven't already):
   ```bash
   git clone https://github.com/deb151292/Remote-Access-Windows-PC.git
   cd Remote-Access-Windows-PC
   ```

2. **Install Dependencies**:
   The application uses Go's standard library and no external Go packages, so no additional dependencies are required.

3. **Run the Application**:
   The application is set to run on port `8080` by default and uses `D:\` as the default base directory.
   ```bash
   go run main.go
   ```
   - This will start the server on `http://localhost:8080`.
   - Open your browser and navigate to `http://localhost:8080` to access the GUI.

4. **Verify the Application**:
   - You should see a folder management interface displaying the contents of the default directory (`D:\`).
   - Features include creating folders, uploading files, searching, and deleting items.

## Changing the Default Path

The default base directory is set to `D:\` in the code. To change this to a different path, you need to modify the `main.go` file. Here's how:

1. **Open the `main.go` File**:
   - Locate the `main.go` file in your project directory.

2. **Find the Default Path**:
   - Look for the line where the `baseDir` variable is defined. It looks like this:
     ```go
     baseDir := "D:\\"
     ```
   - This line is typically near the beginning of the `main` function.

3. **Modify the Path**:
   - Change the value of `baseDir` to your desired directory. For example, to set it to `C:\Users\YourUser\Documents`, update the line to:
     ```go
     baseDir := "C:\\Users\\YourUser\\Documents"
     ```
   - **Note**: Use double backslashes (`\\`) for Windows paths because Go interprets a single backslash as an escape character. For example, `C:\path` should be written as `C:\\path`.

4. **Save the File**:
   - Save your changes to `main.go`.

5. **Restart the Application**:
   - Stop the running application (if it's already running) by pressing `Ctrl+C` in the terminal.
   - Run the application again:
     ```bash
     go run main.go
     ```
   - The application will now use the new default path you specified.

6. **Verify the Change**:
   - Refresh your browser at `http://localhost:8080`.
   - The interface should now display the contents of the new directory you set.

## Exposing the Application with ngrok

To access the application from other devices (e.g., a mobile phone or another computer), you can use ngrok to create a secure public URL for your locally running application. Follow these steps to download, install, and run ngrok, then access your app remotely.

### Step 1: Download and Install ngrok

1. **Create an ngrok Account**:
   - Go to [ngrok.com](https://ngrok.com) and sign up for a free account using your email, Google, or GitHub account.
   - After signing up, you'll be directed to the ngrok dashboard.

2. **Download ngrok**:
   - On the ngrok dashboard, go to the "Getting Started" section or visit the [ngrok download page](https://ngrok.com/download).
   - Select the version for your operating system (e.g., Windows, macOS, or Linux) and download the binary.
   - For Windows, you'll download a `ngrok.zip` file. For macOS/Linux, you'll get a similar zip file or a direct binary.

3. **Extract the Binary**:
   - Unzip the downloaded file to a convenient location on your computer (e.g., `C:\ngrok` on Windows or `~/ngrok` on macOS/Linux).
   - You'll find an executable file named `ngrok` (or `ngrok.exe` on Windows).

4. **Add ngrok to Your System Path** (Optional but Recommended):
   - To run ngrok from any terminal location, add the directory containing `ngrok` to your system's PATH.
   - **Windows**:
     - Right-click on 'This PC' or 'My Computer', select 'Properties', then 'Advanced system settings'.
     - Click 'Environment Variables'.
     - Under 'System Variables' or 'User Variables', find the `Path` variable, click 'Edit', and add the path to the folder containing `ngrok.exe` (e.g., `C:\ngrok`).
   - **macOS/Linux**:
     - Open your terminal and edit your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`):
       ```bash
       echo 'export PATH=$PATH:/path/to/ngrok' >> ~/.bashrc
       source ~/.bashrc
       ```
     - Replace `/path/to/ngrok` with the actual path (e.g., `~/ngrok`).

### Step 2: Authenticate ngrok

1. **Get Your Authtoken**:
   - In the ngrok dashboard, go to the "Your Authtoken" section (usually under "Getting Started").
   - Copy the authtoken provided (it looks like a long string, e.g., `2aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT_xxxxxxxxxxxxxxxxxxxxxx`).

2. **Authenticate ngrok**:
   - Open a terminal and navigate to the directory containing the `ngrok` executable (or just use `ngrok` if you added it to your PATH).
   - Run the following command, replacing `<your_authtoken>` with the token you copied:
     ```bash
     ngrok authtoken <your_authtoken>
     ```
   - You should see a confirmation message like:
     ```
     Authtoken saved to configuration file: /path/to/ngrok.yml
     ```
   - This step links your ngrok client to your account and enables longer tunnel durations and additional features.

### Step 3: Run ngrok to Expose Your Application

1. **Ensure Your Application is Running**:
   - Make sure your Go application is running on `http://localhost:8080` (as described in the "Running the Application" section).

2. **Start an ngrok Tunnel**:
   - In a new terminal window (or the same one if your Go app is running in the background), run:
     ```bash
     ngrok http 8080
     ```
   - ngrok will create a secure tunnel and provide a public URL, such as `https://abc123.ngrok.io`.
   - You'll see output similar to:
     ```
     ngrok by @inconshreveable
     Session Status                online
     Account                       your-account (Plan: Free)
     Version                       3.12.0
     Region                        United States (us)
     Web Interface                 http://127.0.0.1:4040
     Forwarding                    https://abc123.ngrok.io -> http://localhost:8080
     Forwarding                    http://abc123.ngrok.io -> http://localhost:8080
     ```

3. **Access the ngrok Web Interface** (Optional):
   - Open `http://127.0.0.1:4040` in your browser to see the ngrok dashboard.
   - This dashboard shows tunnel status and HTTP request details, which can be useful for debugging.

### Step 4: Access the Application from Other Devices

1. **Share the Public URL**:
   - Copy the `https://abc123.ngrok.io` URL from the ngrok terminal output (use the HTTPS version for security).
   - Share this URL with others or open it on another device (e.g., your phone, tablet, or another computer).

2. **Access the Application**:
   - On the other device, open a browser and navigate to the ngrok URL (e.g., `https://abc123.ngrok.io`).
   - You should see the same folder management GUI as on your local machine.
   - **Note**: The free version of ngrok generates a random URL each time you start a tunnel. If you need a static URL, consider upgrading to a paid ngrok plan and reserving a domain.

3. **Test Features**:
   - Browse folders, create new folders, upload files, and delete items from the remote device.
   - Ensure all features work as expected.

### Step 5: Stop ngrok

- When you're done, stop the ngrok tunnel by pressing `Ctrl+C` in the terminal where ngrok is running.
- This will close the public URL, and your application will no longer be accessible from other devices.

## Security Considerations

- **ngrok Security**:
  - The free version of ngrok creates a publicly accessible URL, meaning anyone with the URL can access your application while the tunnel is active. Be cautious about leaving sensitive features (e.g., file deletion) enabled during testing.
  - Consider using ngrok's paid features like IP restrictions or OAuth to secure your tunnel if needed.

- **Application Security**:
  - The application itself doesn't implement authentication, so anyone with access to the URL can interact with your file system. For production use, add authentication mechanisms to the Go application.

## Troubleshooting

- **Application Not Starting**:
  - Ensure Go is installed correctly (`go version` should return a version number).
  - Check for errors in the terminal when running `go run main.go`. Common issues include invalid paths or permissions.

- **ngrok Tunnel Not Working**:
  - Verify that your application is running on `http://localhost:8080` before starting ngrok.
  - Ensure the port in the ngrok command (`ngrok http 8080`) matches the port your app is using.
  - If you see a "502 Bad Gateway" error when accessing the ngrok URL, your local server might not be running.

- **Access Denied on Other Devices**:
  - Ensure your ngrok tunnel is active and the URL hasn't expired (free tunnels expire after a few hours unless you keep the terminal open).
  - Check if your firewall is blocking inbound connections to port 8080. You may need to allow this port in your firewall settings.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [ngrok Documentation](https://ngrok.com/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
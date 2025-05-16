# Remote File Manager GUI
This is a web-based folder management GUI built with Go and Tailwind CSS. It allows you to browse, create, upload, and delete folders and files on your local system via a browser interface. This README provides instructions on running the application, changing the default path, and exposing it to other devices using LocalTunnel, with support for Windows, Linux, and macOS.

---

## 📸 Preview

![Remote File Manager](screenshot-remote-file-access.png)

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.16 or later): [Download and install Go](https://go.dev/doc/install)
- **Node.js and npm** (required for LocalTunnel): [Download and install Node.js](https://nodejs.org/en/download/)
  - Verify installation by running `node --version` and `npm --version`.
- **Git**: To clone the repository (optional if you already have the code)
- A modern web browser (e.g., Chrome, Firefox)
- An internet connection (for installing LocalTunnel and accessing the app from other devices)

## Running the Application

Follow these steps to run the application locally on your machine (Windows, Linux, or macOS):

1. **Clone the Repository** (if you haven't already):
   ```bash
   git clone https://github.com/deb151292/Remote-Access-PC-free.git
   cd Remote-Access-PC-free
   ```

2. **Install Dependencies**:
   The application uses Go's standard library and no external Go packages, so no additional Go dependencies are required.

3. **Run the Application**:
   The application runs on port `8080` by default. The default base directory depends on your operating system:
   - **Windows**: `D:\`
   - **Linux/macOS**: `~/Documents` (e.g., `/home/user/Documents` on Linux or `/Users/user/Documents` on macOS)

   Run the following command:
   ```bash
   go run main.go
   ```
   - This will start the server on `http://localhost:8080`.
   - Open your browser and navigate to `http://localhost:8080` to access the GUI.

4. **Verify the Application**:
   - You should see a folder management interface displaying the contents of the default directory.
   - Features include creating folders, uploading files, searching, and deleting items.

## Changing the Default Path

The default base directory is automatically set based on your operating system (`D:\` for Windows, `~/Documents` for Linux/macOS). To change this to a different path, you need to modify the `main.go` file. Here's how:

1. **Open the `main.go` File**:
   - Locate the `main.go` file in your project directory.

2. **Find the Default Path**:
   - Look for the `getDefaultBaseDir` function near the beginning of the file. It looks like this:
     ```go
     func getDefaultBaseDir() string {
         switch sysOS := strings.ToLower(os.Getenv("OS")); {
         case strings.Contains(sysOS, "windows"):
             return "D:/" // Windows default
         default:
             homeDir, err := os.UserHomeDir()
             if err != nil {
                 log.Printf("Failed to get user home directory: %v, falling back to /tmp", err)
                 return "/tmp" // Fallback for Linux/macOS
             }
             return filepath.Join(homeDir, "Documents") // e.g., /home/user/Documents or /Users/user/Documents
         }
     }
     ```

3. **Modify the Path**:
   - Update the `return` values in the `getDefaultBaseDir` function to your desired directory.
   - **Windows Example**: To set it to `C:\Users\YourUser\Documents`, modify the Windows case:
     ```go
     return "C:/Users/YourUser/Documents"
     ```
     - **Note**: On Windows, use forward slashes (`/`) or double backslashes (`\\`) in the path string (e.g., `C:\\path` or `C:/path`), as Go interprets a single backslash as an escape character.
   - **Linux/macOS Example**: To set it to `/home/user/MyFolder`, modify the default case:
     ```go
     return "/home/user/MyFolder"
     ```
     - Alternatively, you can use `filepath.Join` to construct the path dynamically:
       ```go
       return filepath.Join("/home/user", "MyFolder")
       ```
     - **Note**: On Linux/macOS, use forward slashes (`/`) and avoid drive letters like `D:`.

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

## Exposing the Application with LocalTunnel

To access the application from other devices (e.g., a mobile phone or another computer), you can use LocalTunnel to create a secure public URL for your locally running application. LocalTunnel adds a layer of security by requiring a password to access the tunnel, which can only be retrieved from the same network. Follow these steps to install and run LocalTunnel, then access your app remotely.

### Step 1: Install LocalTunnel

1. **Ensure Node.js and npm Are Installed**:
   - LocalTunnel is a Node.js package, so you need Node.js and npm installed (see Prerequisites).
   - Verify by running:
     ```bash
     node --version
     npm --version
     ```

2. **Install LocalTunnel Globally**:
   - Open a terminal and install LocalTunnel using npm:
     ```bash
     npm install -g localtunnel
     ```
   - This installs the `lt` command-line tool globally, allowing you to use it from any terminal.

3. **Verify Installation**:
   - Run the following command to check if LocalTunnel is installed correctly:
     ```bash
     lt --version
     ```
   - You should see the version number of LocalTunnel (e.g., `1.9.2`).

### Step 2: Run LocalTunnel to Expose Your Application

1. **Ensure Your Application is Running**:
   - Make sure your Go application is running on `http://localhost:8080` (as described in the "Running the Application" section).

2. **Start a LocalTunnel Session**:
   - In a new terminal window (or the same one if your Go app is running in the background), run:
     ```bash
     lt --port 8080 --print-access-pass
     ```
   - The `--port 8080` flag tells LocalTunnel to tunnel traffic to your application running on port 8080.
   - The `--print-access-pass` flag ensures the access password is displayed in the terminal.

3. **Retrieve the Public URL and Password**:
   - LocalTunnel will generate a public URL, such as `https://<random-subdomain>.loca.lt`.
   - It will also display an access password in the terminal, for example:
     ```
     your url is: https://<random-subdomain>.loca.lt
     Access password: <your-password>
     ```
   - **Where to Find the Password**: The password is displayed in the terminal output on the machine where you ran the `lt` command. It is only visible on the local network where the tunnel is initiated, adding an extra layer of security.
   - **Note**: If you don’t see the password, ensure you included the `--print-access-pass` flag. Without this flag, LocalTunnel may not display the password, and you’d need to check the LocalTunnel dashboard or logs (if available).

4. **Access the LocalTunnel Web Interface** (Optional):
   - LocalTunnel doesn’t provide a built-in web interface like ngrok, but you can visit `https://loca.lt` for more information or to check your tunnel status if you’ve set up a custom subdomain (requires a paid plan).

### Step 3: Access the Application from Other Devices

1. **Share the Public URL**:
   - Copy the URL provided by LocalTunnel (e.g., `https://<random-subdomain>.loca.lt`).
   - Share this URL with others or open it on another device (e.g., your phone, tablet, or another computer).

2. **Enter the Password**:
   - When you or someone else accesses the URL, a password prompt will appear in the browser.
   - Enter the password that was displayed in your terminal (e.g., `<your-password>`).
   - **Security Note**: The password ensures that only those who have access to the local network (where the tunnel was started) can obtain the password and access the application, making LocalTunnel more secure than an open tunnel.
   - Tunnel password can be retrieved from  "curl --location 'https://loca.lt/mytunnelpassword'" Or just paste this > "https://loca.lt/mytunnelpassword" in your browser .

3. **Access the Application**:
   - After entering the correct password, you should see the same folder management GUI as on your local machine.
   - **Note**: LocalTunnel generates a random subdomain each time you start a tunnel. If you need a fixed subdomain, you can use the `--subdomain` flag (e.g., `lt --port 8080 --subdomain mysubdomain `), but this may require a paid plan with LocalTunnel.

4. **Test Features**:
   - Browse folders, create new folders, upload files, and delete items from the remote device.
   - Ensure all features work as expected.

### Step 4: Stop LocalTunnel

- When you're done, stop the LocalTunnel session by pressing `Ctrl+C` in the terminal where LocalTunnel is running.
- This will close the public URL, and your application will no longer be accessible from other devices.

## Security Considerations

- **LocalTunnel Security**:
  - LocalTunnel requires a password to access the tunnel, which adds a layer of security compared to an open tunnel. The password is only displayed on the local network where the tunnel is initiated, ensuring that only authorized users can access it.
  - Be cautious about sharing the password outside your trusted network. If someone else on your local network can see the terminal output, they could retrieve the password.

- **Application Security**:
  - The application itself doesn't implement authentication, so anyone with the URL and password can interact with your file system. For production use, add authentication mechanisms to the Go application to restrict access.

- **File Permissions (Linux/macOS)**:
  - On Linux/macOS, ensure the application has read/write permissions for the directories it accesses. You may need to adjust permissions using `chmod` or run the application with `sudo` if accessing restricted directories.

## Troubleshooting

- **Application Not Starting**:
  - Ensure Go is installed correctly (`go version` should return a version number).
  - Check for errors in the terminal when running `go run main.go`. Common issues include invalid paths or permissions.
  - **Linux/macOS**: If the default path (`~/Documents`) doesn’t exist, create it or update `getDefaultBaseDir` to a valid path:
    ```bash
    mkdir -p ~/Documents
    ```

- **LocalTunnel Not Working**:
  - Verify that your application is running on `http://localhost:8080` before starting LocalTunnel.
  - Ensure the port in the LocalTunnel command (`lt --port 8080`) matches the port your app is using.
  - If you can’t access the URL, double-check that you entered the correct password. The password is case-sensitive.
  - If LocalTunnel fails to start, ensure Node.js and npm are installed correctly, and reinstall LocalTunnel if necessary:
    ```bash
    npm install -g localtunnel
    ```

- **Access Denied on Other Devices**:
  - Ensure your LocalTunnel session is active. The URL expires when you stop the `lt` command.
  - Check if your firewall is blocking inbound connections to port 8080. You may need to allow this port:
    - **Windows**: Use Windows Firewall settings to allow port 8080.
    - **Linux**: Use `ufw` (e.g., `sudo ufw allow 8080`).
    - **macOS**: Use System Preferences > Security & Privacy > Firewall to allow incoming connections.
  - If the password prompt rejects your password, ensure you copied it correctly from the terminal output.

- **Path Issues on Linux/macOS**:
  - If you see errors related to paths, ensure the paths in `main.go` use forward slashes (`/`) and are valid for your system.
  - Check the browser console and server logs for path-related errors.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [LocalTunnel Documentation](https://theboroer.github.io/localtunnel-www/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Node.js Documentation](https://nodejs.org/en/docs/)


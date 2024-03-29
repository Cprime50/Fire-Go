# Fire-Go: Building Modern && Secure Go apps with firebase Auth

## Overview

Fire-Go is an example project that shows how easy it is to use Firebase Auth to build modern Go apps.

   - Go and Gin framework.
   - Sqlite3 for db.
   - Firebase auth with Google Authentication
   - RBAC with firebase
   - JWT is used a session token.
   - OpenApi Docs

## Currently working on
- Writing unit Tests
- Deployment on Aws with Teraform
- Client with nextjs

## Article
This article gives a very detailed guide on this application



## Why Firebase for Authentication in Go?


- **Ease of Integration**: Firebase's Go client library makes it straightforward to integrate Firebase authentication into your Go app .
- **Multiple Authentication Methods**: Firebase supports various authentication methods, including email/password, phone, and social logins, making it versatile for different user bases.
- **Security**: Firebase handles the heavy lifting of security, including token management and encryption, ensuring that your application remains secure.
- **Scalability**: Firebase's infrastructure is designed to scale automatically, making it suitable for applications of all sizes.
- **Real-time Database**: Firebase's real-time database allows for easy data synchronization across clients, enhancing the user experience.



## Getting Started

To get started with this project, follow these steps:

1. **Clone the git repo**: ```git clone https://github.com/Cprime50/Fire-Go```

2. **Set Up Firebase Project**: Create a new Firebase project or use an existing one in your [firebase console](https://console.firebase.google.com)

3. **Configure Authentication**: Set up the desired authentication methods in the Firebase console.

4. **Integrate Firebase Go Client**: Follow the project's [documentation](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&ved=2ahUKEwjnl5XG7feEAxUgTUEAHW3LDbQQFnoECBYQAQ&url=https%3A%2F%2Ffirebaseopensource.com%2Fprojects%2Ffirebase%2Ffirebase-admin-go%2F&usg=AOvVaw1ee2k1xUMEFNFYBKMcoKqU&opi=89978449) to integrate the Firebase Go client into your Go application.

5. **Install the Go dependencys**: cd into project folder and run
```go mod tidy```


6. **Obtain Your Firebase Private Key**:
   - Navigate to the Firebase Console, under project settings, service accounts and download your project's private key.
   - For security, it's recommended to store this key in a `.env` file.

7. **Create .env file**:**Create a `.env` File**:
   - In the root directory of the project, create a `.env` file.
   - Add the following details to the `.env` file:

``` yaml
ADMIN_EMAIL= youremail@mail.com

PORT=:3000

FIREBASE_KEY= your_private_key.json
```

Replace `youremail@mail.com` with your `admin email`, and `path/to/your_private_key.json` with the path to your Firebase private key.

- **Admin Email**: This email will be set as the default admin when authenticated, assuming it doesn't already have an assigned role.


## Contributing

Contributions are welcome! If you have suggestions for improvements or encounter any issues, please feel free to open an issue or submit a pull request.

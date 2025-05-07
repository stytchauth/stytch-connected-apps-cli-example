# Stytch Connected Apps CLI Example

This is a simple CLI application that demonstrates how to authenticate with a Stytch-connected application using OAuth 2.0 with PKCE.

Relevant guide in Stytch docs coming soon!

## Prerequisites

- Go 1.22 or later
- A Stytch project with a Connected App configured
- A Stytch-connected application running at `http://localhost:3000`, and Connected Apps IdentityProvider running at `http://localhost:3000/authorize`.

## Installation

1. Clone this repository:
```bash
git clone https://github.com/stytchauth/stytch-b2b-nextjs-connectedapps-example.git
cd stytch-b2b-nextjs-connectedapps-example
```

2. Install dependencies:
```bash
go mod download
```

## Usage

Run the CLI with your Stytch project credentials:

```bash
go run main.go auth --client-id=YOUR_CONNECTED_APPS_CLIENT_ID --project-id=YOUR_PROJECT_ID
```

The CLI will:
1. Open your default browser to the authorization page
2. Start a local server to receive the OAuth callback
3. Exchange the authorization code for an access token
4. Display the access token and related information

### Command Line Arguments

- `--client-id`: Your Stytch OAuth client ID (required)
- `--project-id`: Your Stytch project ID (required)

## How It Works

1. The CLI generates a PKCE code verifier and challenge
2. Opens your browser to the authorization URL with the PKCE parameters
3. After you authenticate, the web app redirects back to the CLI's local server
4. The CLI exchanges the authorization code for an access token
5. The access token and related information are displayed

## Security

This implementation uses PKCE (Proof Key for Code Exchange) to ensure secure authorization code exchange. The code verifier is generated using cryptographically secure random numbers, and the code challenge is created using SHA-256 hashing.

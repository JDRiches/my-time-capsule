# My Time Capsule
 A little API I made to have fun with Go. Stores a message you create and locks it until the time you specify. Send your future self a message!
# Requirements

## Infrastructure

For the API to connect to the database and authentication, the machine running the API needs to be authenticated with the gcloud CLI.

### Firestore
A Firestore database is required. This database should be named `(default)`

### Firebase Authentication
Authentication is provided by Firebase. A Firebase Authentication provider needs to be set up.

## Environment Variables:
 - `PROJECT_ID`: The project ID of the GCP project which contains the Firestore database


# Running the API

To run the API simply use `go run .`

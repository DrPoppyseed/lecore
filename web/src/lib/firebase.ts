import {
  PUBLIC_FIREBASE_API_KEY,
  PUBLIC_FIREBASE_AUTH_DOMAIN,
  PUBLIC_FIREBASE_PROJECT_ID,
  PUBLIC_FIREBASE_STORAGE_BUCKET,
  PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  PUBLIC_FIREBASE_APP_ID,
  PUBLIC_USE_FIREBASE_AUTH_EMULATOR,
  PUBLIC_FIREBASE_AUTH_EMULATOR_URI,
} from "$env/static/public";
import { initializeApp, type FirebaseApp } from "firebase/app";
import { getAuth, type Auth, connectAuthEmulator } from "firebase/auth";

const config = {
  apiKey: PUBLIC_FIREBASE_API_KEY,
  authDomain: PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: PUBLIC_FIREBASE_APP_ID,
};

let app: FirebaseApp;
let auth: Auth;

export const initializeFirebase = () => {
  if (!app) {
    app = initializeApp(config);
    auth = getAuth(app);

    if (PUBLIC_USE_FIREBASE_AUTH_EMULATOR == "true" || false) {
      connectAuthEmulator(auth, PUBLIC_FIREBASE_AUTH_EMULATOR_URI);
    }
  }
};

export { app, auth };

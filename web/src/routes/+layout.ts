/** @type {import('./$types').LayoutLoad} */

import { initializeFirebase, auth } from "$lib/firebase";
import { browser } from "$app/environment";
import { onAuthStateChanged, type User } from "firebase/auth";

export async function load({ url }) {
  if (browser) {
    try {
      initializeFirebase();
    } catch (ex) {
      console.error(ex);
    }
  }

  function getUser(): Promise<User | null> {
    return new Promise((resolve) => {
      onAuthStateChanged(auth, (user) => resolve(user || null));
    });
  }

  return {
    getUser,
    url: url.pathname,
  };
}

export const ssr = false;

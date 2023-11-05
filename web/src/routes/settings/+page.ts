/** @type {import('./$types').PageLoad} */

import { redirect } from "@sveltejs/kit";

export async function load({ parent }) {
  const { getUser } = await parent();
  const user = await getUser();

  if (!user) {
    throw redirect(302, "/login");
  }
}

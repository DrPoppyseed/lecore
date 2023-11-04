import { z } from "zod";

const PASSWORD_MIN_LEN = 8;

// https://auth0.com/blog/dont-pass-on-the-new-nist-password-guidelines/
// https://bitwarden.com/blog/how-long-should-my-password-be/
// todo: I really want to have some reactive display of how "good" a user's
// password is by showing whether it's been used a bunch, or if it takes a
// very short time to crack them.
// todo: Add a limit to the number of attempts a user can
export const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(PASSWORD_MIN_LEN, {
    message: `Your password must be at least ${PASSWORD_MIN_LEN} characters long`,
  }), // todo: think about how we want to protect against weak passwords.
});

export type LoginSchema = z.infer<typeof loginSchema>;

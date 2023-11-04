<script lang="ts">
  import { Loader2 } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { cn } from "$lib/utils";
  import { goto } from "$app/navigation";
  import { superForm, superValidateSync } from "sveltekit-superforms/client";
  import { signInWithCustomToken } from "firebase/auth";
  import { auth } from "$lib/firebase";
  import { authUser } from "$lib/authStore";
  import { loginSchema, type LoginSchema } from "$lib/schema";
  import { PUBLIC_BACKEND_URI } from "$env/static/public";

  let isLoading: boolean = false;
  let success: boolean | undefined = undefined;
  let className: string | undefined | null = undefined;

  const { form, errors, enhance, constraints, message } = superForm(
    superValidateSync(loginSchema),
    {
      SPA: true,
      validators: loginSchema,
      async onUpdate({ form }) {
        if (form.valid) {
          const token = await createToken(form.data);
          console.log({ token });

          if (token) {
            await loginWithToken(token);
          } else {
            // todo: error handling
            console.log("token missing");
          }
        }
      },
    },
  );

  async function createToken(requestBody: LoginSchema) {
    try {
      const response = await fetch(`${PUBLIC_BACKEND_URI}/login`, {
        method: "POST",
        body: JSON.stringify(requestBody),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const json = (await response.json()) as { token: string };

      if (!response.ok) {
        throw new Error(`${response.status} ${response.statusText}`);
      }

      return json.token;
    } catch (error) {
      // todo: error handling
      console.log(error);
    }
  }

  async function loginWithToken(token: string) {
    try {
      const credentials = await signInWithCustomToken(auth, token);
      $authUser = {
        uid: credentials.user.uid,
        email: credentials.user.email || "",
      };

      goto("/dashboard");
    } catch (error) {
      // todo: error handling
      console.log(error);

      success = false;
    }
  }

  export { className as class };
</script>

{#if $message}
  <p>{$message}</p>
{/if}

<form method="POST" use:enhance>
  <div class={cn("grid gap-6", className)} {...$$restProps}>
    <div class="grid gap-2">
      <Label for="email">Email</Label>
      <Input
        id="email"
        type="email"
        name="email"
        placeholder="name@example.com"
        autocapitalize="none"
        autocorrect="off"
        disabled={isLoading}
        aria-invalid={$errors.email && "true"}
        bind:value={$form.email}
        {...$constraints.email}
      />
    </div>
    <div class="grid gap-2">
      <Label for="password">Password</Label>
      <Input
        id="password"
        type="password"
        name="password"
        autocapitalize="none"
        autocorrect="off"
        disabled={isLoading}
        aria-invalid={$errors.password && "true"}
        bind:value={$form.password}
        {...$constraints.password}
      />
    </div>

    <Button type="submit" value="submit" class="mt-4" disabled={isLoading}>
      {#if isLoading}
        <Loader2 class="mr-2 h-4 w-4 animate-spin" />
      {/if}
      Sign In with Email
    </Button>
  </div>
</form>

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
  import { loginSchema, type LoginSchema } from "$lib/schema";
  import { PUBLIC_BACKEND_URI } from "$env/static/public";

  let isLoading: boolean = false;
  let className: string | undefined | null = undefined;

  const { form, errors, enhance, constraints } = superForm(
    superValidateSync(loginSchema),
    {
      SPA: true,
      validators: loginSchema,
      async onUpdate({ form }) {
        if (form.valid) {
          const token = await register(form.data);

          if (token) {
            await loginWithToken(token);
          } else {
            console.log("token missing");
          }
        }
      },
    },
  );

  async function register(requestBody: LoginSchema) {
    try {
      const response = await fetch(`${PUBLIC_BACKEND_URI}/register`, {
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
      await signInWithCustomToken(auth, token);
      return goto("/dashboard");
    } catch (error) {
      // todo: error handling
      console.log(error);
    }
  }

  export { className as class };
</script>

<form method="POST" use:enhance>
  <div class={cn("grid gap-6", className)} {...$$restProps}>
    <div class="grid gap-2">
      <Label for="register-email">Email</Label>
      <Input
        id="register-email"
        name="email"
        placeholder="name@example.com"
        autocapitalize="none"
        autocomplete="email"
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
        name="password"
        type="password"
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
      Create Account with Email
    </Button>
  </div>
</form>

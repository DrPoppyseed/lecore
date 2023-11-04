<script lang="ts">
  import Button from "$lib/components/ui/button/button.svelte";
  import RegisterForm from "./login-form.svelte";
</script>

<div class="min-h-screen relative flex items-center justify-center">
  <Button
    href="/register"
    variant="ghost"
    class="absolute right-4 top-4 md:right-8 md:top-8"
  >
    Create Account
  </Button>
  <div class="flex w-full mx-8 flex-col justify-center space-y-6 sm:w-[350px]">
    <div class="flex flex-col space-y-2 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">Sign In</h1>
      <p class="text-sm text-muted-foreground">
        Enter your email below to sign in
      </p>
    </div>
    <RegisterForm />
  </div>
</div>

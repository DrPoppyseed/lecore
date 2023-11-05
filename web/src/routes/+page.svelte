<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Brain } from "lucide-svelte";
  import Typewriter from "svelte-typewriter";
  import LL, { locale } from "../i18n/i18n-svelte";
</script>

<div class="relative min-h-screen flex flex-col justify-center items-center">
  <Button
    href="/login"
    variant="ghost"
    class="absolute right-4 top-4 md:right-8 md:top-8"
  >
    Login
  </Button>

  <a href="/dashboard">
    <Brain />
  </a>
  <div class="mt-4 mx-4 flex flex-col items-center">
    <Typewriter
      delay={100}
      wordInterval={3000}
      unwriteInterval={50}
      interval={100}
      mode="loop"
    >
      <h1 class="text-xl md:text-4xl">Leetcode Reviewer</h1>
      <h1 class="text-xl md:text-4xl">Lecore</h1>
    </Typewriter>
    <p class="mt-2">
      {$LL.lecore.desc()}
    </p>
  </div>

  <div class="absolute left-50 bottom-4 md:bottom-8">
    <a href="https://github.com/DrPoppyseed/lecore">Github</a>
  </div>
</div>

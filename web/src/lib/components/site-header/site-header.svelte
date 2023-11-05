<script lang="ts">
  import { Brain } from "lucide-svelte";
  import UserDropdownMenu from "./user-dropdown-menu.svelte";
</script>

<svelte:head>
  <title>lecore|Dashboard</title>
</svelte:head>

<div class="hidden flex-col md:flex">
  <div class="border-b">
    <div class="flex h-16 items-center justify-between px-4 md:px-8">
      <a href="/dashboard">
        <Brain />
      </a>
      <UserDropdownMenu />
    </div>
  </div>
</div>

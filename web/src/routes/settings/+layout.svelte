<script lang="ts">
  import SiteHeader from "$lib/components/site-header/site-header.svelte";
  import { Separator } from "$lib/components/ui/separator";
  import Sidebar from "$lib/components/sidebar.svelte";

  const sidebarNavItems = [
    {
      title: "Account",
      href: "/settings",
    },
  ];
</script>

<SiteHeader />

<div class="hidden space-y-6 p-10 pb-16 md:block">
  <div class="space-y-0.5">
    <h2 class="text-2xl font-bold tracking-tight">Settings</h2>
    <p class="text-muted-foreground">Manage your account settings.</p>
  </div>
  <Separator class="my-6" />
  <div class="flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0">
    <aside class="-mx-4 lg:w-1/5">
      <Sidebar items={sidebarNavItems} />
    </aside>
    <div class="flex-1 lg:max-w-2xl">
      <slot />
    </div>
  </div>
</div>

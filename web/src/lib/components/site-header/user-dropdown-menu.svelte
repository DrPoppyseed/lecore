<script lang="ts">
  import { session } from "$lib/store";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import * as Avatar from "$lib/components/ui/avatar";
  import { Button } from "$lib/components/ui/button";
  import { LogOut, Settings } from "lucide-svelte";
  import { signOut } from "firebase/auth";
  import { auth } from "$lib/firebase";
  import { goto } from "$app/navigation";
  import LL from "../../../i18n/i18n-svelte";

  async function logout() {
    $session.loading = true;
    await signOut(auth)
      .then(() => {
        $session = {
          user: null,
          loading: false,
        };
        goto("/");
      })
      .catch((err) => {
        console.log(err);
        $session.loading = false;
      });
  }
</script>

{#if !$session}
  <div />
{:else}
  <DropdownMenu.Root positioning={{ placement: "bottom-end" }}>
    <DropdownMenu.Trigger asChild let:builder>
      <Button
        variant="ghost"
        builders={[builder]}
        class="relative h-8 w-8 rounded-full"
      >
        <Avatar.Root class="h-8 w-8">
          <Avatar.Image src="/avatars/01.png" alt="fallback" />
          <Avatar.Fallback></Avatar.Fallback>
        </Avatar.Root>
      </Button>
    </DropdownMenu.Trigger>
    <DropdownMenu.Content class="w-56">
      <DropdownMenu.Label class="font-normal">
        <div class="flex flex-col space-y-1">
          <p class="text-sm font-medium leading-none text-muted-foreground">
            {$session.user?.email || $LL.dashboard.anonymous()}
          </p>
        </div>
      </DropdownMenu.Label>
      <DropdownMenu.Separator />
      <DropdownMenu.Group>
        <DropdownMenu.Item>
          <a href="/settings" class="w-full flex items-center">
            <Settings class="opacity-50 h-4 -ml-1 mr-1" />
            {$LL.dashboard.settings()}
          </a>
        </DropdownMenu.Item>
      </DropdownMenu.Group>
      <DropdownMenu.Separator />
      <DropdownMenu.Item
        class="w-full flex items-center cursor-pointer"
        on:click={logout}
      >
        <LogOut class="opacity-50 h-4 -ml-1 mr-1" />
        {$LL.dashboard.logOut()}
      </DropdownMenu.Item>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
{/if}

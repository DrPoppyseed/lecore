<script lang="ts">
  import { onMount } from "svelte";
  import { session } from "$lib/store";
  import { goto } from "$app/navigation";
  import type { LayoutData } from "./$types";
  import { localStorageDetector } from "typesafe-i18n/detectors";
  import { setLocale, locale } from "../i18n/i18n-svelte";
  import type { Locales } from "../i18n/i18n-types";
  import { loadLocaleAsync } from "../i18n/i18n-util.async";
  import { detectLocale } from "../i18n/i18n-util";
  import "../app.postcss";

  export let data: LayoutData;

  onMount(async () => {
    const user = await data.getUser();

    if (!!user) {
      $session.loading = false;
      $session.user = user;
    }

    const wasInAuthRoute = ["/login", "/register"].includes(data.url);
    if ($session.user && wasInAuthRoute) {
      goto("/dashboard");
    }
  });

  onMount(async () => {
    const detectedLocale = detectLocale(localStorageDetector);
    await chooseLocale(detectedLocale);
    localeToSelect = $locale;
  });

  const chooseLocale = async (locale: Locales) => {
    await loadLocaleAsync(locale);
    setLocale(locale);
  };

  let localeToSelect: Locales;
  $: localeToSelect && chooseLocale(localeToSelect);
  $: $locale && localStorage.setItem("lang", $locale);
</script>

<svelte:head>
  <title>lecore</title>
</svelte:head>

<div class="relative min-h-screen">
  <slot />
</div>

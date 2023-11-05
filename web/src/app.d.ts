type Locales = import("$i18n/i18n-types").Locales;
type TranslationFunctions = import("$i18n/i18n-types").TranslationFun;

declare global {
  namespace App {
    interface Locals {
      locale: Locales;
      LL: TranslationFunctions;
    }
  }
}

export {};

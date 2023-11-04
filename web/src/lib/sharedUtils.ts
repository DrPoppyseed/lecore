import type { z } from "zod";

export function formDataToJSON<Schema extends z.ZodTypeAny>(
  formData: FormData,
  schema: Schema,
) {
  // https://github.com/colinhacks/zod/discussions/967#discussioncomment-2244720
  return schema.parse(Object.fromEntries(formData)) as z.infer<Schema>;
}

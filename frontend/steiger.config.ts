import fsd from "@feature-sliced/steiger-plugin";
import { defineConfig } from "steiger";

export default defineConfig([
  ...fsd.configs.recommended,
  {
    files: ["./src/shared/api/**", "./src/shared/hooks/**"],
    rules: {
      "fsd/segments-by-purpose": "off",
      "fsd/no-reserved-folder-names": "off",
    },
  },
]);

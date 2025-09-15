module.exports = {
  admin: {
    output: {
      mode: "tags-split",
      target: "src/shared/api/admin",
      schemas: "src/shared/api/admin/model",
      client: "swr",
      override: {
        mutator: {
          path: "./src/shared/api/client.ts",
          name: "adminApiClient",
        },
      },
      biome: true,
      clean: true,
    },
    input: {
      target: "../infrastructure/schemas/openapi/admin/openapi.yaml",
    },
  },
};

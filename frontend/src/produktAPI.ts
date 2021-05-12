import * as z from "zod";

const DataProduktTilgangSchema = z.object({
  subject: z.string(),
  start: z.string(),
  end: z.string(),
});

const BucketStoreSchema = z.object({
  type: z.literal("bucket"),
  project_id: z.string(),
  bucket_id: z.string(),
});

const BigQuerySchema = z.object({
  type: z.literal("bigquery"),
  project_id: z.string(),
  dataset_id: z.string(),
  resource_id: z.string(),
});

export const DataLagerSchema = z.union([BucketStoreSchema, BigQuerySchema]);
export const DataProduktSchema = z.object({
  name: z.string(),
  description: z.string().nullable(),
  owner: z.string(),
  datastore: DataLagerSchema.array().nullable(),
  access: DataProduktTilgangSchema.array(),
});

export const DataProduktResponseSchema = z.object({
  id: z.string(),
  updated: z.string(),
  created: z.string(),
  data_product: DataProduktSchema,
});

export const BrukerInfoSchema = z.object({
  email: z.string(),
  teams: z.array(z.string()),
});

const DataProduktListSchema = DataProduktResponseSchema.array();

export type DataProdukt = z.infer<typeof DataProduktSchema>;
export type DataProduktTilgang = z.infer<typeof DataProduktTilgangSchema>;
export type DataProduktResponse = z.infer<typeof DataProduktResponseSchema>;
export type DataProduktListe = z.infer<typeof DataProduktListSchema>;
export type DataLager = z.infer<typeof DataLagerSchema>;
export type BrukerInfo = z.infer<typeof BrukerInfoSchema>;

export const hentProdukter = async (): Promise<DataProduktListe> => {
  let apiURL = "http://localhost:8080/api/v1/dataproducts";
  const res = await fetch(apiURL);
  const json = await res.json();

  return DataProduktListSchema.parse(json);
};

export const hentBrukerInfo = async (): Promise<BrukerInfo> => {
  let apiURL = "http://localhost:8080/api/v1/userinfo";
  const res = await fetch(apiURL, { credentials: "include" });
  const json = await res.json();

  return BrukerInfoSchema.parse(json);
};

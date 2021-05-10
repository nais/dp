import * as z from "zod";

const DataProduktTilgangSchema = z.object({
  subject: z.string(),
  start: z.string(),
  end: z.string(),
});

const DataLagerSchema = z.object({
  project_id: z.string(),
  dataset_id: z.string(),
  type: z.string(),
});

export const DataProduktSchema = z
  .object({
    name: z.string(),
    description: z.string().nullable(),
    owner: z.string(),
    datastore: DataLagerSchema.partial(),
    access: DataProduktTilgangSchema.array(),
  })
  .partial();

export const DataProduktResponseSchema = z.object({
  id: z.string(),
  updated: z.string(),
  created: z.string(),
  data_product: DataProduktSchema,
});

const DataProduktListSchema = DataProduktResponseSchema.array();

export type DataProdukt = z.infer<typeof DataProduktSchema>;
export type DataProduktResponse = z.infer<typeof DataProduktResponseSchema>;
export type DataProduktListe = z.infer<typeof DataProduktListSchema>;

export const hentProdukter = async (): Promise<DataProduktListe> => {
  let apiURL = "http://localhost:8080/api/v1/dataproducts";
  const res = await fetch(apiURL);
  const json = await res.json();

  return DataProduktListSchema.parse(json);
};

import * as z from "zod";

const DataProduktTilgangSchema = z.object({
  subject: z.string(),
  start: z.string(),
  end: z.string(),
});

const DataProduktSchema = z.object({
  id: z.string(),
  updated: z.string(),
  created: z.string(),
  data_product: z.object({
    name: z.string(),
    description: z.string(),
    owner: z.string(),
    access: DataProduktTilgangSchema.array(),
    uri: z.string(),
  }),
});
const DataProduktListSchema = DataProduktSchema.array();

export type DataProdukt = z.infer<typeof DataProduktSchema>;
export type DataProduktListe = z.infer<typeof DataProduktListSchema>;

export const hentProdukter = async (): Promise<DataProduktListe> => {
  let apiURL = "http://localhost:8080/api/v1/dataproducts";
  const res = await fetch(apiURL);
  const json = await res.json();

  return DataProduktListSchema.parse(json);
};

import * as z from "zod";
export const BACKEND_ENDPOINT =
  process.env.BACKEND_ENDPOINT || "http://localhost:8080";
export const API_ROOT = `${BACKEND_ENDPOINT}/api/v1`;

const DataProduktTilgangSchema = z.object({
  subject: z.string(),
  expires: z.string().nullable(),
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

const DataProduktPartialSchema = DataProduktSchema.deepPartial();

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
export type DataProduktPartial = z.infer<typeof DataProduktPartialSchema>;
export type DataProduktTilgang = z.infer<typeof DataProduktTilgangSchema>;
export type DataProduktResponse = z.infer<typeof DataProduktResponseSchema>;
export type DataProduktListe = z.infer<typeof DataProduktListSchema>;
export type DataLager = z.infer<typeof DataLagerSchema>;
export type BrukerInfo = z.infer<typeof BrukerInfoSchema>;

export const hentProdukter = async (): Promise<DataProduktListe> => {
  const res = await fetch(`${API_ROOT}/dataproducts`);
  const json = await res.json();

  return DataProduktListSchema.parse(json);
};

export const hentProdukt = async (
  produktID: string
): Promise<DataProduktResponse> => {
  try {
    const res = await fetch(`${API_ROOT}/dataproducts/${produktID}`, {
      credentials: "include",
    });

    if (res.status !== 200) {
      throw new Error(`Feil: ${await res.text()}`);
    } else {
      return DataProduktResponseSchema.parse(await res.json());
    }
  } catch (e) {
    console.log(e);
    throw new Error(`Nettverksfeil: ${e}`);
  }
};

export const slettProdukt = async (produktID: string): Promise<void> => {
  try {
    const res = await fetch(`${API_ROOT}/dataproducts/${produktID}`, {
      method: "delete",
      credentials: "include",
    });

    if (res.status !== 204) {
      throw new Error(`Feil: ${await res.text()}`);
    } else {
      return;
    }
  } catch (e) {
    console.log(e);
    throw new Error(`Nettverksfeil: ${e}`);
  }
};

export const opprettProdukt = async (
  nyttProdukt: DataProdukt
): Promise<string> => {
  const res = await fetch(`${API_ROOT}/dataproducts`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(nyttProdukt),
  });

  if (res.status !== 201) {
    throw new Error(
      `Kunne ikke opprette nytt produkt: ${res.status}: ${await res.text()}`
    );
  }

  const newID = await res.text();
  return newID;
};

export const oppdaterProdukt = async (
  produktID: string,
  oppdatertProdukt: DataProduktPartial
): Promise<string> => {
  const res = await fetch(`${API_ROOT}/dataproducts/${produktID}`, {
    method: "PUT",
    credentials: "include",
    body: JSON.stringify(oppdatertProdukt),
  });

  if (res.status !== 204) {
    throw new Error(
      `Kunne ikke oppdatere produkt: ${res.status}: ${await res.text()}`
    );
  }

  const newID = await res.text();
  return newID;
};

export const giTilgang = async (
  produkt: DataProduktResponse,
  subject: string,
  expiry: Date | null
) => {
  const expiryString = expiry instanceof Date ? expiry.toISOString() : null;

  const produktOppdateringer: DataProduktPartial = {
    access: [
      {
        subject: subject,
        expires: expiryString,
      },
    ],
  };

  await oppdaterProdukt(produkt.id, produktOppdateringer);
};

export const hentBrukerInfo = async (): Promise<BrukerInfo> => {
  const res = await fetch(`${API_ROOT}/userinfo`, { credentials: "include" });
  const json = await res.json();

  // dummy values, please replace later
  let user = BrukerInfoSchema.parse(json);
  return user;
};

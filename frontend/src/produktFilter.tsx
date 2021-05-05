import React, { useEffect, useState } from "react";
import { Select } from "nav-frontend-skjema";
import { ProduktListeState } from "./produktListe";

interface ProduktFilterProps {
  state: ProduktListeState;
  dispatch: any; // chickening out again
}

export const ProduktFilter = ({
  state,
  dispatch,
}: ProduktFilterProps): JSX.Element => {
  const [productOwners, setProductOwners] = useState<string[]>([]);

  useEffect(() => {
    const productOwners: string[] = [];

    state.products.forEach((x) => {
      if (
        x?.data_product?.owner &&
        !productOwners.includes(x.data_product.owner)
      ) {
        productOwners.push(x.data_product.owner);
      }
      setProductOwners(productOwners);
    });
  }, [state.products]);

  return (
    <Select label={"Filtrer på produkteier"}>
      {productOwners.map((o) => (
        <option value={o}>{o}</option>
      ))}
    </Select>
  );
};

export default ProduktFilter;

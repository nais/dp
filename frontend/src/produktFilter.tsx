import React, { ChangeEvent, useEffect, useState } from "react";
import { Select } from "nav-frontend-skjema";
import { ProduktListeState } from "./produktListe";
import "./produktFilter.less";

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

  const selectTeam = (e: ChangeEvent<HTMLSelectElement>) => {
    dispatch({
      type: "FILTER_CHANGE",
      filter: e.target.value,
    });
  };

  return (
    <Select label={"Filtrer på produkteier"} onChange={(e) => selectTeam(e)}>
      <option key="xxx" value="">
        {"Velg team"}
      </option>
      {productOwners.map((o) => (
        <option key={o} value={o}>
          {o}
        </option>
      ))}
    </Select>
  );
};

export default ProduktFilter;

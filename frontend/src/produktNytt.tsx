import React from "react";

import { opprettProdukt, DataProdukt } from "./produktAPI";
import "react-datepicker/dist/react-datepicker.css";
import { useHistory } from "react-router-dom";
import ProduktSkjema from "./produktSkjema";

export const ProduktNytt = (): JSX.Element => {
  const history = useHistory();

  const handleNewProduct = (produkt: DataProdukt): void => {
    opprettProdukt(produkt).then((newID) => history.push(`/produkt/${newID}`));
  };

  return (
    <div style={{ margin: "1em 1em 0 1em" }}>
      <ProduktSkjema onProductReady={handleNewProduct} />
    </div>
  );
};

export default ProduktNytt;

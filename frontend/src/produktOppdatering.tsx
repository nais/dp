import {useHistory, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {DataProdukt, DataProduktResponse, hentProdukt, oppdaterProdukt, opprettProdukt} from "./produktAPI";
import ProduktSkjema from "./produktSkjema";

export const ProduktOppdatering = (): JSX.Element => {
    const history = useHistory();

    let { produktID } = useParams<{ produktID: string}>();
    const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
    const [error, setError] = useState<string | null>();

    const handleProduct = (produkt: DataProdukt): void => {
        oppdaterProdukt(produktID, produkt).then(() => history.push(`/produkt/${produktID}`));
    }

    useEffect(() => {
        hentProdukt(produktID)
            .then((p) => {
                setProdukt(p);
                setError(null);
            })
            .catch((e) => {
                setError(e.toString());
            });
    }, [produktID]);

    return (
        <div style={{ margin: "1em 1em 0 1em" }}>
            { produkt && <ProduktSkjema produkt={produkt.data_product}
                                        avbryt={()=>history.push(`/produkt/${produktID}`)}
                                        onProductReady={handleProduct}/>
            }
        </div>

    )
};

export default ProduktOppdatering
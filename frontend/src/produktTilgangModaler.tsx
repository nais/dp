import { Systemtittel } from "nav-frontend-typografi";
import React, { useState } from "react";
import { DataProduktResponse, slettProdukt } from "./produktAPI";
import Modal from "nav-frontend-modal";
import { useHistory } from "react-router-dom";
import { Fareknapp } from "nav-frontend-knapper";

interface SlettProduktProps {
  produktID: string;
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

interface GiTilgangProps {
  produkt: DataProduktResponse;
  tilgangIsOpen: boolean;
  setTilgangIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

export const GiTilgang = ({
  produkt,
  tilgangIsOpen,
  setTilgangIsOpen,
}: GiTilgangProps): JSX.Element => {
  return (
    <Modal
      isOpen={tilgangIsOpen}
      onRequestClose={() => setTilgangIsOpen(false)}
      closeButton={true}
      contentLabel="Gi tilgang"
    >
      <div className="slette-bekreftelse">
        <Systemtittel>Gi tilgang</Systemtittel>
      </div>
    </Modal>
  );
};

export const SlettProdukt = ({
  produktID,
  isOpen,
  setIsOpen,
}: SlettProduktProps): JSX.Element => {
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const deleteProduct = async (id: string) => {
    try {
      await slettProdukt(id);
      history.push("/");
    } catch (e) {
      setError(e.toString());
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onRequestClose={() => setIsOpen(false)}
      closeButton={true}
      contentLabel="Min modalrute"
    >
      <div className="slette-bekreftelse">
        <Systemtittel>Er du sikker?</Systemtittel>
        {error ? <p>{error}</p> : null}
        <Fareknapp onClick={() => deleteProduct(produktID)}>
          {error ? "Prøv igjen" : "Ja"}
        </Fareknapp>
      </div>
    </Modal>
  );
};

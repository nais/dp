import { Systemtittel } from "nav-frontend-typografi";
import React, { useState } from "react";
import { DataProduktResponse, slettProdukt } from "./produktAPI";
import Modal from "nav-frontend-modal";
import { useHistory } from "react-router-dom";
import { Fareknapp } from "nav-frontend-knapper";
import { Select, SkjemaGruppe } from "nav-frontend-skjema";

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
      appElement={document.getElementById("app") || undefined}
      isOpen={tilgangIsOpen}
      onRequestClose={() => setTilgangIsOpen(false)}
      closeButton={true}
      contentLabel="Gi tilgang"
    >
      <div className="gi-tilgang">
        <Systemtittel>Gi tilgang</Systemtittel>

        <SkjemaGruppe>
          <Select>
            <option value="til meg">til meg</option>
            <option value="til team">til team</option>
          </Select>
        </SkjemaGruppe>
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
      appElement={document.getElementById("app") || undefined}
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

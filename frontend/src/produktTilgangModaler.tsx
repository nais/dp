import { Feilmelding, Systemtittel } from "nav-frontend-typografi";
import React, { useState, useContext } from "react";
import { giTilgang, DataProduktResponse, slettProdukt } from "./produktAPI";
import { UserContext } from "./userContext";
import Modal from "nav-frontend-modal";
import { ToggleGruppe } from "nav-frontend-toggle";
import { useHistory } from "react-router-dom";
import { Fareknapp, Hovedknapp } from "nav-frontend-knapper";

import DatePicker from "react-datepicker";
import "./produktTilgangModaler.less";

interface SlettProduktProps {
  produktID: string;
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

export const GiTilgang: React.FC<{
  produkt: DataProduktResponse;
  tilgangIsOpen: boolean;
  callback: (hasChanged: boolean) => void;
}> = ({ produkt, tilgangIsOpen, callback }) => {
  const [endDate, setEndDate] = useState<Date | null>(new Date());
  const [evig, setEvig] = useState<boolean>(false);
  const [feilmelding, setFeilmelding] = useState<string | null>(null);
  const userContext = useContext(UserContext);
  if (!userContext) return null;

  const handleSubmit = async () => {
    try {
      await giTilgang(produkt, userContext.email, endDate);
      setFeilmelding(null);
      callback(true);
    } catch (e) {
      setFeilmelding(e.toString());
    }
  };

  return (
    <Modal
      appElement={document.getElementById("app") || undefined}
      isOpen={tilgangIsOpen}
      onRequestClose={() => callback(false)}
      closeButton={false}
      contentLabel="Gi tilgang"
      className={"gitilgang"}
    >
      <Systemtittel>Gi tilgang til {userContext.email}?</Systemtittel>
      {feilmelding ? <Feilmelding>{feilmelding}</Feilmelding> : null}
      <ToggleGruppe
        minstEn={true}
        defaultToggles={[
          {
            children: "velg sluttdato...",
            pressed: true,
            onClick: (e) => setEvig(false),
          },
          { children: "evig", onClick: (e) => setEvig(true) },
        ]}
      />

      {!evig ? (
        <div className={"datovalg"}>
          <DatePicker
            selected={endDate}
            onChange={(e) => setEndDate(e as Date)}
            selectsEnd
            endDate={endDate}
            startDate={new Date()}
            minDate={new Date()}
            inline
          />
        </div>
      ) : null}
      <div className={"knapperad"}>
        <Fareknapp onClick={() => callback(false)}>Avbryt</Fareknapp>
        <Hovedknapp className={"bekreft"} onClick={() => handleSubmit()}>
          Bekreft
        </Hovedknapp>
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

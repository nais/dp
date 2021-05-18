import { Systemtittel } from "nav-frontend-typografi";
import React, { useState, useContext } from "react";
import {
  giTilgang,
  DataProduktResponse,
  slettProdukt,
  DataProdukt,
} from "./produktAPI";
import { UserContext } from "./userContext";
import Modal from "nav-frontend-modal";
import { ToggleKnapp } from "nav-frontend-toggle";
import { useHistory } from "react-router-dom";
import { Fareknapp, Hovedknapp } from "nav-frontend-knapper";
import { Select, SkjemaGruppe } from "nav-frontend-skjema";
import DatePicker from "react-datepicker";
import "./produktTilgangModaler.less";

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

export const GiTilgang: React.FC<{
  produkt: DataProduktResponse;
  tilgangIsOpen: boolean;
  setTilgangIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}> = ({ produkt, tilgangIsOpen, setTilgangIsOpen }) => {
  const [startDate, setStartDate] = useState<Date>(new Date());
  const [endDate, setEndDate] = useState<Date | null>(null);
  const [evig, setEvig] = useState<boolean>(false);

  const userContext = useContext(UserContext);
  if (!userContext) return null;

  const handleDateUpdate = (newValue: Date | [Date, Date] | null) => {
    console.log("new value:", typeof newValue, newValue);
    if (newValue instanceof Date) setEndDate(newValue);
  };

  const handleSubmit = () => {
    //giTilgang(produkt, userContext.email, endDate, endDate)
  };

  return (
    <Modal
      appElement={document.getElementById("app") || undefined}
      isOpen={tilgangIsOpen}
      onRequestClose={() => setTilgangIsOpen(false)}
      closeButton={true}
      contentLabel="Gi tilgang"
    >
      <div className="gitilgang">
        <Systemtittel>Gi tilgang til {userContext.email}?</Systemtittel>
        <ToggleKnapp onClick={() => setEvig(!evig)}>Evig tilgang</ToggleKnapp>
        {!evig ? (
          <DatePicker
            selected={endDate}
            onChange={(e) => handleDateUpdate(e)}
            selectsEnd
            endDate={endDate}
            startDate={startDate}
            minDate={startDate}
          />
        ) : null}
        <Hovedknapp onClick={() => {}}>Ja</Hovedknapp>
        <Fareknapp onClick={() => setTilgangIsOpen(false)}>Nei</Fareknapp>
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

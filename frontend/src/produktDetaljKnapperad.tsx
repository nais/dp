import { Fareknapp, Knapp } from "nav-frontend-knapper";
import React, { useContext, useState } from "react";
import { UserContext } from "./userContext";
import {
  DataProduktResponse,
  DataProduktTilgangListe,
  getCurrentAccessState,
  hentProdukt,
  isOwner,
} from "./produktAPI";

export const ProduktKnapperad: React.FC<{
  produkt: DataProduktResponse;
  tilganger: DataProduktTilgangListe;
  openSlett: () => void;
  openTilgang: () => void;
}> = ({ produkt, tilganger, openSlett, openTilgang }) => {
  const userContext = useContext(UserContext);

  const harTilgang = (tilganger: DataProduktTilgangListe): boolean => {
    if (!tilganger) return false;
    const tilgangerBehandlet = getCurrentAccessState(tilganger);
    if (!tilgangerBehandlet) return false;

    for (const tilgang of tilgangerBehandlet) {
      if (tilgang.subject == userContext?.email) {
        if (tilgang?.expires && new Date(tilgang.expires) > new Date()) {
          return true;
        }
      }
    }

    return false;
  };

  return (
    <div className="knapperad">
      {isOwner(produkt?.data_product, userContext?.teams) && (
        <Fareknapp onClick={() => openSlett()}>Slett</Fareknapp>
      )}

      {userContext && !harTilgang(tilganger) && (
        <Knapp onClick={() => openTilgang()}>Få tilgang</Knapp>
      )}
    </div>
  );
};

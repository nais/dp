import React, {useContext, useEffect, useState} from "react";
import {
    DataProduktResponse,
    DataProduktTilgangListe,
    DataProduktTilgangResponse,
    getCurrentAccessState,
    deleteAccess,
    isOwner,
} from "./produktAPI";
import {UserContext} from "./userContext";
import "moment/locale/nb";
import moment from "moment";
import {
    Feilmelding,
    Undertekst,
    Undertittel,
} from "nav-frontend-typografi";
import {Child} from "@navikt/ds-icons";
import "./produktDetaljTilganger.less";
import {Xknapp} from "nav-frontend-ikonknapper";

export const ProduktTilganger: React.FC<{
    produkt: DataProduktResponse | null;
    tilganger: DataProduktTilgangListe | null;
    refreshAccessState: () => void
}> = ({produkt, tilganger,refreshAccessState}) => {
    const userContext = useContext(UserContext);
    const [error, setError] = useState<string | null>(null);

    const entryShouldBeDisplayed = (subject: string | undefined): boolean => {
        if (!produkt?.data_product || !userContext?.teams) return false;
        // Hvis produkteier, vis all tilgang;
        if (isOwner(produkt?.data_product, userContext?.teams)) return true;
        // Ellers, vis kun dine egne tilganger.
        return subject === userContext?.email;
    };

    if (!tilganger) return <></>;

    const tilgangsLinje = (tilgang: DataProduktTilgangResponse) => {
        const accessEnd = moment(tilgang.expires).format("LLL");

        return (
            <>
                <Child style={{height: "100%", width: "auto"}}/>
                <div className={"tilgangsTekst"}>
                    <Undertittel>{tilgang.subject}</Undertittel>
                    <Undertekst>
                        Innvilget av <em>{tilgang.author}</em> til {accessEnd}
                    </Undertekst>
                </div>
                <Xknapp
                    mini
                    type={"fare"}
                    onClick={async () => {
                        if (produkt?.id && tilgang?.subject)
                            try {
                                await deleteAccess(produkt.id, tilgang.subject, "user");
                                refreshAccessState()
                            } catch (e) {
                                setError(e.toString())
                            }
                    }}
                />
            </>
        );
    };

    const synligeTilganger = getCurrentAccessState(
        tilganger
            .filter((tilgang) => tilgang.action !== "verify")
            .filter((tilgang) => entryShouldBeDisplayed(tilgang.subject))
    );

    if (!synligeTilganger?.length)
        return <p>Ingen relevante tilganger definert</p>;

    return (
        <div className={"produkt-tilganger"}>
            {error && <Feilmelding>{error}</Feilmelding>}

            {synligeTilganger.map((tilgang, index) => (
                <div key={index} className={"produkt-tilgang"}>
                    {tilgangsLinje(tilgang)}
                </div>
            ))}
        </div>
    );
};

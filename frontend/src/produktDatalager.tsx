import React from "react";
import { DataLager } from "./produktAPI";
import { BigQueryIcon } from "./svgIcons";
import "./produktDatalager.less";
import { Undertittel } from "nav-frontend-typografi";

export const DatalagerInfo: React.FC<{ ds: DataLager }> = ({ ds }) => {
  const BigQueryEntry = (
    <div className={"bigqueryentry datalagerentry"}>
      <BigQueryIcon />
      <Undertittel>BigQuery</Undertittel>
    </div>
  );
  const BucketEntry = <></>;

  return BigQueryEntry;
};

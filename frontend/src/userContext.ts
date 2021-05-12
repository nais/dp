import { createContext } from "react";
import { API_ROOT, BrukerInfo, BrukerInfoSchema } from "./produktAPI";

export const UserContext = createContext<BrukerInfo | null>(null);

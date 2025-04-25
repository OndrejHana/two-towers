import { getAuth } from "../auth";
import { navigate } from "../routes";

export function renderRoot() {
  const auth = getAuth();
  if (!auth.session) {
    navigate("/login");
  } else {
    navigate("/main");
  }
}

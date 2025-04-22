import { renderLoginPage } from "./auth";
import { renderError } from "./components/error";
import { renderLoading } from "./components/loading";
import { renderLobby } from "./components/lobby";
import { renderMain } from "./components/main";
import { renderGame } from "./components/game";

import { Clerk } from "@clerk/clerk-js";

export const LOADING = 0;
export const LOGIN = 1;
export const MAIN = 2;
export const LOBBY = 3;
export const GAME = 4;
export const ERROR = 5;

export class State {
  /**
   * @param {Element} appDiv
   * @param {number} [initialState=LOADING]
   */
  constructor(appDiv, initialState = LOADING) {
    this.appDiv = appDiv;
    this.state = initialState;
    this.context = {};
    this.error = null;
    this.controller = null;
  }

  /**
   * @param {number} state
   */
  renderWith(state) {
    if (state === this.state) {
      console.log(state, this.state);
      return;
    }
    this.state = state;
    return this.render();
  }

  /**
   * @returns {Error | null}
   */
  render() {
    if (this.controller) {
      this.controller.abort();
    }

    switch (this.state) {
      case LOADING:
        renderLoading(this.appDiv);
        break;
      case LOGIN:
        if (
          !this.context.hasOwnProperty("clerk") ||
          (!this.context.clerk) instanceof Clerk
        ) {
          return new Error("Clerk not in context");
        } else {
          renderLoginPage(this.appDiv, this.context.clerk);
        }
        break;
      case MAIN:
        if (
          this.context.hasOwnProperty("clerk") &&
          this.context.clerk instanceof Clerk
        ) {
          this.controller = renderMain(this.appDiv, this.context.clerk, this);
        } else {
          return new Error("Clerk not in context");
        }
        break;
      case LOBBY:
        if (
          !this.context.hasOwnProperty("clerk") ||
          (!this.context.clerk) instanceof Clerk
        ) {
          return new Error("Clerk not in context");
        }

        renderLobby(this.appDiv, this.context.clerk, this);
        break;
      case GAME:
        if (!this.context.hasOwnProperty("payload")) {
          return new Error("Payload not in context");
        }

        renderGame(this.appDiv, this.context.payload, null);
        break;
      default:
        renderError(this.appDiv, this.error);
    }

    return null;
  }
}

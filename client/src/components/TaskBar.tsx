import * as React from "react";

import { classNames } from "../utils";
import { GamePlayer } from "../objects/State";
import { Game } from "../objects/Game";

export type OptionsProps = {
    player?: GamePlayer;
    game: Game;
    sendAction: (action: any) => void;
}

export class TaskBar extends React.Component<OptionsProps, {}> {
    public getOptions(): string[] {
        if (!this.props.player) {
            return [];
        }

        let options: string[] = Object.keys(this.props.player.options || {});

        if (this.props.player.id == this.props.game.owner_id) {
            options.push('Start');
        }

        return options;
    }

    public startGameAction() {
        let action = {
            type: 'admin',
            data: {
                action: 'start',
            },
        };

        this.props.sendAction(action);
    }

    public sendGameAction(type: string) {
        // TODO: get action type
        // TODO: get action data
        //this.props.sendAction(action);
    }

    public render() {
        console.log("Player for task bar:", this.props.player);
        const keys = this.props.player && this.props.player.options ? Object.keys(this.props.player.options) : [];

        return !this.props.player ? '' : (
            <div className={classNames("table-task-bar indigo darken-4", {"active": keys.length > 0})}>
                <div className="row">
                    {this.getOptions().map((key) => {
                        return <div>
                            <label>
                                {key}:
                                <button onClick={this.startGameAction.bind(this)} type="button">{key}</button>
                            </label>
                        </div>;
                    })}
                </div>
            </div>
        );
    }
}
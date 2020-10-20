import * as React from "react";

import { classNames } from "../../utils";
import { Game } from "../../objects/Game";
import { EventState } from "../../objects/State";

type ChipSettingProps = {
    es: EventState;
    game: Game;
    playerID: string;
    requests: any[];
    sendAction: (type: string, data: any) => void;
}

type ChipSettingsState = {
    chipRequest: number;
}

export class ChipSettings extends React.Component<ChipSettingProps, ChipSettingsState> {
    constructor(props: any) {
        super(props);

        this.state = {chipRequest: this.getDefaultBuyIn()};
    }

    private setChipsRequest(event: any) {
        this.setState({chipRequest: event.target.value});
    }

    private sendChipsRequest(event: any) {
        this.props.sendAction(
            'admin',
            {
                action: 'chip_request',
                value: this.state.chipRequest
            }
        )
    }

    public sendChipsResponse(id: string) {
        this.props.sendAction('admin', { action: 'chip_response', value: id })
    }

    private getDefaultBuyIn(): number {
        return Math.max(
            Math.min(
                this.props.game.options.buy_in_min * 10,
                this.props.game.options.buy_in_max
            ),
            this.props.game.options.buy_in_min,
        )
    }

    private hasChips(): boolean {
        let player = this.props.es.getPlayer(this.props.playerID)

        return player && player.stack > 0;
    }

    public render() {
        return <tbody>
            <tr>
                <th colSpan={4}>Chip Requests</th>
            </tr>

            {this.props.requests.map((req) => {
                return <tr>
                    <td colSpan={2}>{this.props.es.getPlayer(req.player_id).username}</td>
                    <td>
                        <button onClick={this.sendChipsResponse.bind(this, req.player_id.toString())} className="btn-flat green lighten-2" type="button">
                            Approve
                        </button>
                    </td>
                    <td>
                        <button onClick={this.sendChipsResponse.bind(this, '-' + req.player_id.toString())} className="btn-flat red lighten-2" type="button">
                            Deny
                        </button>
                    </td>
                </tr>
            })}
            <tr className={classNames({'setting-highlight': !this.hasChips()})}>
                <td colSpan={2}>Request Chips</td>
                <td>
                    <input type="number" step={this.props.game.options.big_blind}
                        defaultValue={this.getDefaultBuyIn()}
                        min={this.props.game.options.buy_in_min} max={this.props.game.options.buy_in_max}
                        name="chip_request" onChange={this.setChipsRequest.bind(this)} />
                </td>
                <td>
                    <button onClick={this.sendChipsRequest.bind(this)} type="button" className="btn-flat">
                        Submit
                    </button>
                </td>
            </tr>
            <tr></tr>
        </tbody>;
    }
}

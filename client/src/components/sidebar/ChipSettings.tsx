import * as React from "react";

import { classNames } from "../../utils";
import { Game } from "../../objects/Game";
import { AnonymousPlayer, AnonymousPlayerWithChips, findPlayer } from "../../objects/Player";

type ChipSettingProps = {
    game: Game;
    player: AnonymousPlayerWithChips;
    players: AnonymousPlayer[];
    requests: any[];
    sendAction: (type: string, data: any) => void;
}

type ChipSettingsState = {
    chipRequest: number;
    disableRequestButton: boolean;
}

export class ChipSettings extends React.Component<ChipSettingProps, ChipSettingsState> {
    constructor(props: any) {
        super(props);

        this.state = {
            chipRequest: this.getDefaultBuyIn(),
            disableRequestButton: false,
        };
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
        );

        // disable button temporarily
        this.setState({ disableRequestButton: true });
        setTimeout(
            this.setState.bind(this, { disableRequestButton: false }),
            3000,
        );
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
        return this.props.player && this.props.player.stack > 0;
    }

    public render() {
        return <>
            <tr>
                <th colSpan={4}>Chip Requests</th>
            </tr>

            {this.props.requests.map((req, i) => {
                return <tr key={i}>
                    <td colSpan={2}>{findPlayer(req.player_id, this.props.players)?.username}</td>
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
        </>;
    }
}

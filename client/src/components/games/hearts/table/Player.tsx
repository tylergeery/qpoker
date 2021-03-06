import * as React from "react";
import { GamePlayer, Manager } from "../objects/State";
import { ConnectionHandler } from "../../../../connection/ws";
import { Card } from "../../../../objects/Card";
import { Game } from "../../../../objects/Game";
import { createGameAction, classNames } from "../../../../utils";

type PlayerProps = {
    playerID: number;
    player: GamePlayer;
    index: number;
    manager: Manager;
    game?: Game;
    cards: Card[];
    conn: ConnectionHandler;
}

type PlayerState = {
    timer: number;
}

type HandProps = {
    gameState: string;
    cards: Card[];
}

class Hand extends React.Component<HandProps, {}> {
    render() {
        return this.props.gameState != 'Init' ? (
            <div>
                <img className="card" src={`/assets/media/cards/${this.props.cards[0].imageName()}.svg`} />
                <img className="card" src={`/assets/media/cards/${this.props.cards[1].imageName()}.svg`} />
            </div>
        ) : '';
    }
}

type HandTimerProps = {
    allIn: boolean;
    timer: number;
    decisionTime: number;
}

class HandTimer extends React.Component<HandTimerProps, {}> {
    private getTimerWidth(): string {
        return (100 / this.props.decisionTime * this.props.timer).toString() + "%";
    }

    render() {
        if (!this.props.timer) {
            return '';
        }

        return <div className="progress">
            <div className="determinate flip-x" style={{width: this.getTimerWidth()}}></div>
        </div>
    }
}


class HandActions extends React.Component<PlayerProps, {}> {
    bet: number;

    constructor(props: any) {
        super(props)
        this.bet = this.defaultBetAmount();
    }

    public getOptions(): string[] {
        let options: string[] = [];
        let ordered: string[] = ['can_fold', 'can_check', 'can_call', 'can_bet'].reverse();

        if (!this.props.player.options) {
            return options
        }

        for (let i = 0; i < ordered.length; i++) {
            if (this.props.player.options[ordered[i]]) {
                options.push(ordered[i].slice(4));  // Remove `can_` prefix
            }
        }

        return options;
    }

    public sendAction(event: React.ChangeEvent<HTMLInputElement>) {
        switch (event.target.innerHTML) {
            case 'bet':
                this.props.conn.send(createGameAction({
                    name: event.target.innerHTML,
                    amount: this.bet || this.defaultBetAmount(),
                }));
                break;
            default:
                this.props.conn.send(createGameAction({
                    name: event.target.innerHTML,
                    amount: 0,
                }));
                break;
        }
    }

    private registerBet(event: React.ChangeEvent<HTMLInputElement>) {
        this.bet = parseInt(event.target.value);
    }

    private defaultBetAmount(): number {
        let bigBlind = this.props.game.options.big_blind;

        return Math.min(
            this.props.player.stack,
            bigBlind + this.props.manager.pot.toCover(this.props.player.id)
        );
    }

    render() {
        let options = this.getOptions();
        let betAmount = this.defaultBetAmount();

        return this.props.playerID.toString() === this.props.player.id.toString() ? (
            <div className={classNames('global-action-bar', {'active': !!options.length})}>
                {(options.length && options[0] === 'bet') ? (
                    <div className="input-field inline">
                        <input type="number" defaultValue={betAmount}
                            onChange={this.registerBet.bind(this)} />
                    </div>
                ) : ''}
                {options.map((opt) =>
                    <button key={opt}
                            onClick={this.sendAction.bind(this)} type="button"
                            className="waves-light btn-small">
                        {opt
                    }</button>
                )}
            </div>
        ) : '';
    }
}

export class PlayerSpotlight extends React.Component<PlayerProps, PlayerState> {
    protected isWinner(): boolean {
        return this.props.manager.pot.payouts[+this.props.player.id] > 0;
    }

    render() {
        return <div className={classNames("spotlight", {active: this.isWinner()})}></div>
    }
}

export class Player extends React.Component<PlayerProps, PlayerState> {
    interval: number = null

    constructor(props: any) {
        super(props)
        this.interval = null;
        this.state = { timer: 0 };
    }

    protected getCountdownTime(): number {
        let ts = +(new Date()) / 1000
        let seconds = this.props.game?.options.decision_time - (ts - this.props.manager.state.table.activeAt);

        if (!seconds || seconds < 0) {
            seconds = 0;
        }

        return Math.floor(seconds);
    }

    protected isSelected(): boolean {
        if (this.props.manager.pot.payouts && Object.keys(this.props.manager.pot.payouts).length) {
            return false;
        }

        if (this.props.manager.allIn) {
            return false;
        }

        if (this.getCountdownTime() <= 0) {
            return false
        }

        return this.props.manager.state.table.active === this.props.player.id;
    }

    private countDown(seconds: number) {
        if (seconds <= 0 || !this.isSelected() || this.props.manager.status === 'paused') {
            this.interval = null;
            this.setState({timer: 0});
            return;
        }

        this.setState({timer: seconds});
        window.setTimeout(this.countDown.bind(this, seconds-1), 1000);
    }

    private startTimer(): any {
        if (this.interval || !this.props.manager.state.table.activeAt) {
            return
        }
        this.interval = 1;

        this.countDown(this.getCountdownTime());
    }

    render() {
        if (this.isSelected()) {
            this.startTimer()
        }

        return <div className={ `player table-pos-${this.props.index}` }>
            <video id={`player-video-${this.props.player.id}`} autoPlay></video>
            <small>{`${this.props.player.username} (${this.props.player.stack})` }</small>
            <Hand gameState={this.props.manager.state.state} cards={this.props.cards} />
            <HandActions {...this.props} />
            <HandTimer allIn={this.props.manager.allIn} timer={this.state.timer} decisionTime={this.props.game?.options.decision_time} />
            <p>{this.props.manager.pot.playerBets[+this.props.player.id] || ''}</p>
            <PlayerSpotlight {...this.props} />
        </div>
    }
}

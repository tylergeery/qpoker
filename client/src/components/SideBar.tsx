import * as React from "react";

import { classNames } from "../utils";
import { EventState } from "../objects/State";
import { Game } from "../objects/Game";
import { Settings } from "./sidebar/Settings";
import { Chat } from "./sidebar/Chat";
import { History } from "./sidebar/History";
import { ConnectionHandler } from "../connection/ws";

type SideBarProps = {
    es: EventState;
    game?: Game;
    playerID: string;
    playerToken: string;
    conn: ConnectionHandler;
}

type SideBarState = {
    selectedTab: string;
}

export class SideBar extends React.Component<SideBarProps, SideBarState> {
    constructor(props: any) {
        super(props);
        this.state = { selectedTab: "settings" };
    }

    public getNavOptions(): string[] {
        return ['history', 'chat', 'settings'];
    }

    public getNavLabel(option: string): string {
        if (option != 'settings') {
            return option;
        }

        if (this.props.game.owner_id.toString() != this.props.playerID.toString()) {
            return option;
        }

        return 'admin';
    }

    public navSelect(event: any) {
        this.setState({ selectedTab: event.target.innerHTML });

        event.stopPropagation();
    }

    public render() {
        return <div className="sidebar grey-text">
            <div className="sidebar-nav">
                {this.getNavOptions().map((option) => {
                    let label = this.getNavLabel(option);

                    return <a href="#"
                                className={classNames("btn-flat", { underline: this.state.selectedTab === label })}
                                onClick={this.navSelect.bind(this)} >
                        {label}
                    </a>;
                })}
            </div>
            {this.props.game ? (
                <div>
                    <Chat {...this.props} active={this.state.selectedTab === 'chat'}/>
                    <History {...this.props} active={this.state.selectedTab === 'history'}/>
                    <Settings {...this.props} active={this.state.selectedTab === 'settings' || this.state.selectedTab === 'admin'}/>
                </div>
            ) : ''}
        </div>;
    }
}
/* @flow */
import React, {Component} from 'react'
import {Box} from '../common-adapters'
import SearchHelp from './help.desktop'
import {globalStyles} from '../styles/style-guide'
import UserSearch from './user-search/render'
import type {Props} from './render'

class Render extends Component<void, Props, void> {
  _renderComingSoon () {
    return <SearchHelp username={this.props.username} />
  }

  _renderInfoPane () {
    // TODO(mm) maybe shell out to a smart component here?
    // One that only takes a username and can figure out the rest
    return <Box style={{flex: 1}} />
  }

  render () {
    if (this.props.showComingSoon) {
      return this._renderComingSoon()
    }

    return (
      <Box style={globalStyles.flexBoxRow}>
        <UserSearch
          searchHintText={this.props.searchHintText}
          searchText={this.props.searchText}
          searchIcon={this.props.searchIcon}
          selectedService={this.props.selectedService}
          results={this.props.results}
          onClickService={this.props.onClickService}
          onClickResult={this.props.onClickResult}
          onSearch={this.props.onSearch} />
        {this._renderInfoPane()}
      </Box>
    )
  }
}

export default Render

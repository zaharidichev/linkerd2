import Paper from '@material-ui/core/Paper';
import PropTypes from 'prop-types';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import _ from 'lodash';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing.unit * 3,
    marginBottom: theme.spacing.unit * 3,
    overflowX: 'auto',
  },
  table: {},
  activeSortIcon: {
    opacity: 1,
  },
  inactiveSortIcon: {
    opacity: 0.4,
  },
  denseTable: {
    paddingRight: "8px",
    paddingLeft: "8px"
  }
});

class BaseTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      order: this.props.defaultOrder || "asc",
      orderBy: this.props.defaultOrderBy
    };
  }

  createSortHandler = col => () => {
    let orderBy = col.dataIndex;
    let order = col.defaultSortOrder || 'asc';

    if (this.state.orderBy === orderBy && this.state.order === order) {
      order = order === 'asc' ? 'desc' : 'asc';
    }

    this.setState({ order, orderBy });
  };

  sortRows = (tableRows, tableColumns, order, orderBy) => {
    if (!orderBy) {
      return tableRows;
    }

    let col = _.find(tableColumns, ['dataIndex', orderBy]);
    let sorted = tableRows.sort(col.sorter);
    return order === 'desc' ? _.reverse(sorted) : sorted;
  }

  renderHeaderCell = (col, order, orderBy) => {
    let active = orderBy === col.dataIndex;
    const { classes, padding } = this.props;
    let tableCell;

    if (col.sorter) {
      tableCell = (
        <TableCell
          key={col.key || col.dataIndex}
          numeric={col.isNumeric}
          sortDirection={orderBy === col.dataIndex ? order : false}
          className={classNames({[classes.denseTable]: padding === 'dense'})}>
          <TableSortLabel
            active={active}
            direction={active ? order : col.defaultSortOrder || 'asc'}
            classes={{icon: active ? classes.activeSortIcon : classes.inactiveSortIcon}}
            onClick={this.createSortHandler(col)}>
            {_.isEmpty(col.tooltip) ? col.title : <Tooltip title={col.tooltip} placement="top"><React.Fragment>{col.title}</React.Fragment></Tooltip>}
            {/* <Tooltip title={col.tooltip} placement="top"><React.Fragment>{col.title}</React.Fragment></Tooltip> */}
          </TableSortLabel>
        </TableCell>
      );
    } else {
      tableCell = (
        <TableCell
          key={col.key || col.dataIndex}
          numeric={col.isNumeric}
          className={classNames({[classes.denseTable]: padding === 'dense'})}>
          {col.title}
        </TableCell>
      );
    }

    return _.isNil(col.tooltip) ? tableCell : <Tooltip key={col.key || col.dataIndex} placement="top" title="FOOBAR">{tableCell}</Tooltip>;
  }

  render() {
    const { classes, tableRows, tableColumns, tableClassName, rowKey, padding} = this.props;
    const {order, orderBy} = this.state;
    const sortedTableRows = this.sortRows(tableRows, tableColumns, order, orderBy);

    return (
      <Paper className={classes.root}>
        <Table className={`${classes.table} ${tableClassName}`} padding={padding}>
          <TableHead>
            <TableRow>
              { _.map(tableColumns, c => (
                this.renderHeaderCell(c, order, orderBy)
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {
              _.map(sortedTableRows, d => {
              let key = !rowKey ? d.key : rowKey(d);
              return (
                <TableRow key={key}>
                  { _.map(tableColumns, c => (
                    <TableCell
                      className={classNames({[classes.denseTable]: padding === 'dense'})}
                      key={`table-${key}-${c.key || c.dataIndex}`}
                      numeric={c.isNumeric}>
                      {c.render ? c.render(d) : _.get(d, c.dataIndex)}
                    </TableCell>
                  ))}
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Paper>
    );
  }
}

BaseTable.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  defaultOrder: PropTypes.string,
  defaultOrderBy: PropTypes.string,
  padding: PropTypes.string,
  rowKey: PropTypes.func,
  tableClassName: PropTypes.string,
  tableColumns: PropTypes.arrayOf(PropTypes.shape({
    dataIndex: PropTypes.string,
    defaultSortOrder: PropTypes.string,
    isNumeric: PropTypes.bool,
    render: PropTypes.func,
    sorter: PropTypes.func,
    title: PropTypes.string
  })).isRequired,
  tableRows: PropTypes.arrayOf(PropTypes.shape({}))
};

BaseTable.defaultProps = {
  defaultOrder: "asc",
  defaultOrderBy: null,
  padding: "default",
  rowKey: null,
  tableClassName: "",
  tableRows: [],
};

export default withStyles(styles)(BaseTable);

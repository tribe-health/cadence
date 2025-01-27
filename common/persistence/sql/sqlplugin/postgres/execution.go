// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package postgres

import (
	"context"
	"database/sql"

	"github.com/uber/cadence/common/persistence/sql/sqlplugin"
)

const (
	executionsColumns = `shard_id, domain_id, workflow_id, run_id, next_event_id, last_write_version, data, data_encoding`

	createExecutionQuery = `INSERT INTO executions(` + executionsColumns + `)
 VALUES(:shard_id, :domain_id, :workflow_id, :run_id, :next_event_id, :last_write_version, :data, :data_encoding)`

	updateExecutionQuery = `UPDATE executions SET
 next_event_id = :next_event_id, last_write_version = :last_write_version, data = :data, data_encoding = :data_encoding
 WHERE shard_id = :shard_id AND domain_id = :domain_id AND workflow_id = :workflow_id AND run_id = :run_id`

	getExecutionQuery = `SELECT ` + executionsColumns + ` FROM executions
 WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4`

	listExecutionQuery = `SELECT ` + executionsColumns + ` FROM executions
 WHERE shard_id = $1 AND workflow_id > $2 ORDER BY workflow_id LIMIT $3`

	deleteExecutionQuery = `DELETE FROM executions
 WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4`

	lockExecutionQueryBase = `SELECT next_event_id FROM executions
 WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4`

	writeLockExecutionQuery = lockExecutionQueryBase + ` FOR UPDATE`
	readLockExecutionQuery  = lockExecutionQueryBase + ` FOR SHARE`

	createCurrentExecutionQuery = `INSERT INTO current_executions
(shard_id, domain_id, workflow_id, run_id, create_request_id, state, close_status, start_version, last_write_version) VALUES
(:shard_id, :domain_id, :workflow_id, :run_id, :create_request_id, :state, :close_status, :start_version, :last_write_version)`

	deleteCurrentExecutionQuery = "DELETE FROM current_executions WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4"

	getCurrentExecutionQuery = `SELECT
shard_id, domain_id, workflow_id, run_id, create_request_id, state, close_status, start_version, last_write_version
FROM current_executions WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3`

	lockCurrentExecutionJoinExecutionsQuery = `SELECT
ce.shard_id, ce.domain_id, ce.workflow_id, ce.run_id, ce.create_request_id, ce.state, ce.close_status, ce.start_version, e.last_write_version
FROM current_executions ce
INNER JOIN executions e ON e.shard_id = ce.shard_id AND e.domain_id = ce.domain_id AND e.workflow_id = ce.workflow_id AND e.run_id = ce.run_id
WHERE ce.shard_id = $1 AND ce.domain_id = $2 AND ce.workflow_id = $3 FOR UPDATE`

	lockCurrentExecutionQuery = getCurrentExecutionQuery + ` FOR UPDATE`

	updateCurrentExecutionsQuery = `UPDATE current_executions SET
run_id = :run_id,
create_request_id = :create_request_id,
state = :state,
close_status = :close_status,
start_version = :start_version,
last_write_version = :last_write_version
WHERE
shard_id = :shard_id AND
domain_id = :domain_id AND
workflow_id = :workflow_id
`

	getTransferTasksQuery = `SELECT task_id, data, data_encoding
 FROM transfer_tasks WHERE shard_id = $1 AND task_id > $2 AND task_id <= $3 ORDER BY shard_id, task_id LIMIT $4`

	createTransferTasksQuery = `INSERT INTO transfer_tasks(shard_id, task_id, data, data_encoding)
 VALUES(:shard_id, :task_id, :data, :data_encoding)`

	deleteTransferTaskQuery             = `DELETE FROM transfer_tasks WHERE shard_id = $1 AND task_id = $2`
	rangeDeleteTransferTaskQuery        = `DELETE FROM transfer_tasks WHERE shard_id = $1 AND task_id > $2 AND task_id <= $3`
	rangeDeleteTransferTaskByBatchQuery = `DELETE FROM transfer_tasks WHERE shard_id = $1 AND task_id IN (SELECT task_id FROM
		transfer_tasks WHERE shard_id = $1 AND task_id > $2 AND task_id <= $3 ORDER BY task_id LIMIT $4)`

	getCrossClusterTasksQuery = `SELECT task_id, data, data_encoding
 FROM cross_cluster_tasks WHERE target_cluster = $1 AND shard_id = $2 AND task_id > $3 AND task_id <= $4 ORDER BY task_id LIMIT $5`

	createCrossClusterTasksQuery = `INSERT INTO cross_cluster_tasks(target_cluster, shard_id, task_id, data, data_encoding)
 VALUES(:target_cluster, :shard_id, :task_id, :data, :data_encoding)`

	deleteCrossClusterTaskQuery             = `DELETE FROM cross_cluster_tasks WHERE target_cluster = $1 AND shard_id = $2 AND task_id = $3`
	rangeDeleteCrossClusterTaskQuery        = `DELETE FROM cross_cluster_tasks WHERE target_cluster = $1 AND shard_id = $2 AND task_id > $3 AND task_id <= $4`
	rangeDeleteCrossClusterTaskByBatchQuery = `DELETE FROM cross_cluster_tasks WHERE target_cluster = $1 AND shard_id = $2 AND task_id IN (SELECT task_id FROM
		cross_cluster_tasks WHERE target_cluster = $1 AND shard_id = $2 AND task_id > $3 AND task_id <= $4 ORDER BY task_id LIMIT $5)`

	createTimerTasksQuery = `INSERT INTO timer_tasks (shard_id, visibility_timestamp, task_id, data, data_encoding)
  VALUES (:shard_id, :visibility_timestamp, :task_id, :data, :data_encoding)`

	getTimerTasksQuery = `SELECT visibility_timestamp, task_id, data, data_encoding FROM timer_tasks
  WHERE shard_id = $1
  AND ((visibility_timestamp >= $2 AND task_id >= $3) OR visibility_timestamp > $4)
  AND visibility_timestamp < $5
  ORDER BY visibility_timestamp,task_id LIMIT $6`

	deleteTimerTaskQuery             = `DELETE FROM timer_tasks WHERE shard_id = $1 AND visibility_timestamp = $2 AND task_id = $3`
	rangeDeleteTimerTaskQuery        = `DELETE FROM timer_tasks WHERE shard_id = $1 AND visibility_timestamp >= $2 AND visibility_timestamp < $3`
	rangeDeleteTimerTaskByBatchQuery = `DELETE FROM timer_tasks WHERE shard_id = $1 AND (visibility_timestamp,task_id) IN (SELECT visibility_timestamp,task_id FROM
		timer_tasks WHERE shard_id = $1 AND visibility_timestamp >= $2 AND visibility_timestamp < $3 ORDER BY visibility_timestamp,task_id LIMIT $4)`

	createReplicationTasksQuery = `INSERT INTO replication_tasks (shard_id, task_id, data, data_encoding)
  VALUES(:shard_id, :task_id, :data, :data_encoding)`

	getReplicationTasksQuery = `SELECT task_id, data, data_encoding FROM replication_tasks WHERE
shard_id = $1 AND
task_id > $2 AND
task_id <= $3
ORDER BY task_id LIMIT $4`

	deleteReplicationTaskQuery             = `DELETE FROM replication_tasks WHERE shard_id = $1 AND task_id = $2`
	rangeDeleteReplicationTaskQuery        = `DELETE FROM replication_tasks WHERE shard_id = $1 AND task_id <= $2`
	rangeDeleteReplicationTaskByBatchQuery = `DELETE FROM replication_tasks WHERE shard_id = $1 AND task_id IN (SELECT task_id FROM
		replication_tasks WHERE task_id <= $2 ORDER BY task_id LIMIT $3)`

	getReplicationTasksDLQQuery = `SELECT task_id, data, data_encoding FROM replication_tasks_dlq WHERE
source_cluster_name = $1 AND
shard_id = $2 AND
task_id > $3 AND
task_id <= $4
ORDER BY task_id LIMIT $5`
	getReplicationTaskDLQQuery = `SELECT count(1) as count FROM replication_tasks_dlq WHERE
source_cluster_name = $1 AND
shard_id = $2`

	bufferedEventsColumns     = `shard_id, domain_id, workflow_id, run_id, data, data_encoding`
	createBufferedEventsQuery = `INSERT INTO buffered_events(` + bufferedEventsColumns + `)
VALUES (:shard_id, :domain_id, :workflow_id, :run_id, :data, :data_encoding)`

	deleteBufferedEventsQuery = `DELETE FROM buffered_events WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4`
	getBufferedEventsQuery    = `SELECT data, data_encoding FROM buffered_events WHERE shard_id = $1 AND domain_id = $2 AND workflow_id = $3 AND run_id = $4`

	insertReplicationTaskDLQQuery = `
INSERT INTO replication_tasks_dlq
            (source_cluster_name,
             shard_id,
             task_id,
             data,
             data_encoding)
VALUES     (:source_cluster_name,
            :shard_id,
            :task_id,
            :data,
            :data_encoding)
`
	deleteReplicationTaskFromDLQQuery = `
	DELETE FROM replication_tasks_dlq
		WHERE source_cluster_name = $1
		AND shard_id = $2
		AND task_id = $3`

	rangeDeleteReplicationTaskFromDLQQuery = `
	DELETE FROM replication_tasks_dlq
		WHERE source_cluster_name = $1
		AND shard_id = $2
		AND task_id > $3
		AND task_id <= $4`
	rangeDeleteReplicationTaskFromDLQByBatchQuery = `DELETE FROM replication_tasks_dlq WHERE source_cluster_name = $1 AND shard_id = $2 AND task_id IN (SELECT task_id FROM
		replication_tasks_dlq WHERE source_cluster_name = $1 AND shard_id = $2 AND task_id > $3 AND task_id <= $4 ORDER BY task_id LIMIT $5)`
)

// InsertIntoExecutions inserts a row into executions table
func (pdb *db) InsertIntoExecutions(ctx context.Context, row *sqlplugin.ExecutionsRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createExecutionQuery, row)
}

// UpdateExecutions updates a single row in executions table
func (pdb *db) UpdateExecutions(ctx context.Context, row *sqlplugin.ExecutionsRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, updateExecutionQuery, row)
}

// SelectFromExecutions reads a single row from executions table
// The list execution query result is order by workflow ID only. It may returns duplicate record with pagination.
func (pdb *db) SelectFromExecutions(ctx context.Context, filter *sqlplugin.ExecutionsFilter) ([]sqlplugin.ExecutionsRow, error) {
	var rows []sqlplugin.ExecutionsRow
	var err error
	if len(filter.DomainID) == 0 && filter.Size > 0 {
		err = pdb.conn.SelectContext(ctx, &rows, listExecutionQuery, filter.ShardID, filter.WorkflowID, filter.Size)
		if err != nil {
			return nil, err
		}
	} else {
		var row sqlplugin.ExecutionsRow
		err = pdb.conn.GetContext(ctx, &row, getExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, err
}

// DeleteFromExecutions deletes a single row from executions table
func (pdb *db) DeleteFromExecutions(ctx context.Context, filter *sqlplugin.ExecutionsFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
}

// ReadLockExecutions acquires a write lock on a single row in executions table
func (pdb *db) ReadLockExecutions(ctx context.Context, filter *sqlplugin.ExecutionsFilter) (int, error) {
	var nextEventID int
	err := pdb.conn.GetContext(ctx, &nextEventID, readLockExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
	return nextEventID, err
}

// WriteLockExecutions acquires a write lock on a single row in executions table
func (pdb *db) WriteLockExecutions(ctx context.Context, filter *sqlplugin.ExecutionsFilter) (int, error) {
	var nextEventID int
	err := pdb.conn.GetContext(ctx, &nextEventID, writeLockExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
	return nextEventID, err
}

// InsertIntoCurrentExecutions inserts a single row into current_executions table
func (pdb *db) InsertIntoCurrentExecutions(ctx context.Context, row *sqlplugin.CurrentExecutionsRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createCurrentExecutionQuery, row)
}

// UpdateCurrentExecutions updates a single row in current_executions table
func (pdb *db) UpdateCurrentExecutions(ctx context.Context, row *sqlplugin.CurrentExecutionsRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, updateCurrentExecutionsQuery, row)
}

// SelectFromCurrentExecutions reads one or more rows from current_executions table
func (pdb *db) SelectFromCurrentExecutions(ctx context.Context, filter *sqlplugin.CurrentExecutionsFilter) (*sqlplugin.CurrentExecutionsRow, error) {
	var row sqlplugin.CurrentExecutionsRow
	err := pdb.conn.GetContext(ctx, &row, getCurrentExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID)
	return &row, err
}

// DeleteFromCurrentExecutions deletes a single row in current_executions table
func (pdb *db) DeleteFromCurrentExecutions(ctx context.Context, filter *sqlplugin.CurrentExecutionsFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteCurrentExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
}

// LockCurrentExecutions acquires a write lock on a single row in current_executions table
func (pdb *db) LockCurrentExecutions(ctx context.Context, filter *sqlplugin.CurrentExecutionsFilter) (*sqlplugin.CurrentExecutionsRow, error) {
	var row sqlplugin.CurrentExecutionsRow
	err := pdb.conn.GetContext(ctx, &row, lockCurrentExecutionQuery, filter.ShardID, filter.DomainID, filter.WorkflowID)
	return &row, err
}

// LockCurrentExecutionsJoinExecutions joins a row in current_executions with executions table and acquires a
// write lock on the result
func (pdb *db) LockCurrentExecutionsJoinExecutions(ctx context.Context, filter *sqlplugin.CurrentExecutionsFilter) ([]sqlplugin.CurrentExecutionsRow, error) {
	var rows []sqlplugin.CurrentExecutionsRow
	err := pdb.conn.SelectContext(ctx, &rows, lockCurrentExecutionJoinExecutionsQuery, filter.ShardID, filter.DomainID, filter.WorkflowID)
	return rows, err
}

// InsertIntoTransferTasks inserts one or more rows into transfer_tasks table
func (pdb *db) InsertIntoTransferTasks(ctx context.Context, rows []sqlplugin.TransferTasksRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createTransferTasksQuery, rows)
}

// SelectFromTransferTasks reads one or more rows from transfer_tasks table
func (pdb *db) SelectFromTransferTasks(ctx context.Context, filter *sqlplugin.TransferTasksFilter) ([]sqlplugin.TransferTasksRow, error) {
	var rows []sqlplugin.TransferTasksRow
	err := pdb.conn.SelectContext(ctx, &rows, getTransferTasksQuery, filter.ShardID, filter.MinTaskID, filter.MaxTaskID, filter.PageSize)
	if err != nil {
		return nil, err
	}
	return rows, err
}

// DeleteFromTransferTasks deletes one or more rows from transfer_tasks table
func (pdb *db) DeleteFromTransferTasks(ctx context.Context, filter *sqlplugin.TransferTasksFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteTransferTaskQuery, filter.ShardID, filter.TaskID)
}

// RangeDeleteFromTransferTasks deletes multi rows from transfer_tasks table
func (pdb *db) RangeDeleteFromTransferTasks(ctx context.Context, filter *sqlplugin.TransferTasksFilter) (sql.Result, error) {
	if filter.PageSize > 0 {
		return pdb.conn.ExecContext(ctx, rangeDeleteTransferTaskByBatchQuery, filter.ShardID, filter.MinTaskID, filter.MaxTaskID, filter.PageSize)
	}
	return pdb.conn.ExecContext(ctx, rangeDeleteTransferTaskQuery, filter.ShardID, filter.MinTaskID, filter.MaxTaskID)
}

// InsertIntoCrossClusterTasks inserts one or more rows into cross_cluster_tasks table
func (pdb *db) InsertIntoCrossClusterTasks(ctx context.Context, rows []sqlplugin.CrossClusterTasksRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createCrossClusterTasksQuery, rows)
}

// SelectFromCrossClusterTasks reads one or more rows from cross_cluster_tasks table
func (pdb *db) SelectFromCrossClusterTasks(ctx context.Context, filter *sqlplugin.CrossClusterTasksFilter) ([]sqlplugin.CrossClusterTasksRow, error) {
	var rows []sqlplugin.CrossClusterTasksRow
	err := pdb.conn.SelectContext(ctx, &rows, getCrossClusterTasksQuery, filter.TargetCluster, filter.ShardID, filter.MinTaskID, filter.MaxTaskID, filter.PageSize)
	if err != nil {
		return nil, err
	}
	return rows, err
}

// DeleteFromCrossClusterTasks deletes one or more rows from cross_cluster_tasks table
func (pdb *db) DeleteFromCrossClusterTasks(ctx context.Context, filter *sqlplugin.CrossClusterTasksFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteCrossClusterTaskQuery, filter.TargetCluster, filter.ShardID, filter.TaskID)
}

// RangeDeleteFromCrossClusterTasks deletes multi rows from cross_cluster_tasks table
func (pdb *db) RangeDeleteFromCrossClusterTasks(ctx context.Context, filter *sqlplugin.CrossClusterTasksFilter) (sql.Result, error) {
	if filter.PageSize > 0 {
		return pdb.conn.ExecContext(ctx, rangeDeleteCrossClusterTaskByBatchQuery, filter.TargetCluster, filter.ShardID, filter.MinTaskID, filter.MaxTaskID, filter.PageSize)
	}
	return pdb.conn.ExecContext(ctx, rangeDeleteCrossClusterTaskQuery, filter.TargetCluster, filter.ShardID, filter.MinTaskID, filter.MaxTaskID)
}

// InsertIntoTimerTasks inserts one or more rows into timer_tasks table
func (pdb *db) InsertIntoTimerTasks(ctx context.Context, rows []sqlplugin.TimerTasksRow) (sql.Result, error) {
	for i := range rows {
		rows[i].VisibilityTimestamp = pdb.converter.ToPostgresDateTime(rows[i].VisibilityTimestamp)
	}
	return pdb.conn.NamedExecContext(ctx, createTimerTasksQuery, rows)
}

// SelectFromTimerTasks reads one or more rows from timer_tasks table
func (pdb *db) SelectFromTimerTasks(ctx context.Context, filter *sqlplugin.TimerTasksFilter) ([]sqlplugin.TimerTasksRow, error) {
	var rows []sqlplugin.TimerTasksRow
	filter.MinVisibilityTimestamp = pdb.converter.ToPostgresDateTime(filter.MinVisibilityTimestamp)
	filter.MaxVisibilityTimestamp = pdb.converter.ToPostgresDateTime(filter.MaxVisibilityTimestamp)
	err := pdb.conn.SelectContext(ctx, &rows, getTimerTasksQuery, filter.ShardID, filter.MinVisibilityTimestamp,
		filter.TaskID, filter.MinVisibilityTimestamp, filter.MaxVisibilityTimestamp, filter.PageSize)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].VisibilityTimestamp = pdb.converter.FromPostgresDateTime(rows[i].VisibilityTimestamp)
	}
	return rows, err
}

// DeleteFromTimerTasks deletes one or more rows from timer_tasks table
func (pdb *db) DeleteFromTimerTasks(ctx context.Context, filter *sqlplugin.TimerTasksFilter) (sql.Result, error) {
	filter.VisibilityTimestamp = pdb.converter.ToPostgresDateTime(filter.VisibilityTimestamp)
	return pdb.conn.ExecContext(ctx, deleteTimerTaskQuery, filter.ShardID, filter.VisibilityTimestamp, filter.TaskID)
}

// RangeDeleteFromTimerTasks deletes multi rows from timer_tasks table
func (pdb *db) RangeDeleteFromTimerTasks(ctx context.Context, filter *sqlplugin.TimerTasksFilter) (sql.Result, error) {
	filter.MinVisibilityTimestamp = pdb.converter.ToPostgresDateTime(filter.MinVisibilityTimestamp)
	filter.MaxVisibilityTimestamp = pdb.converter.ToPostgresDateTime(filter.MaxVisibilityTimestamp)
	if filter.PageSize > 0 {
		return pdb.conn.ExecContext(ctx, rangeDeleteTimerTaskByBatchQuery, filter.ShardID, filter.MinVisibilityTimestamp, filter.MaxVisibilityTimestamp, filter.PageSize)
	}
	return pdb.conn.ExecContext(ctx, rangeDeleteTimerTaskQuery, filter.ShardID, filter.MinVisibilityTimestamp, filter.MaxVisibilityTimestamp)
}

// InsertIntoBufferedEvents inserts one or more rows into buffered_events table
func (pdb *db) InsertIntoBufferedEvents(ctx context.Context, rows []sqlplugin.BufferedEventsRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createBufferedEventsQuery, rows)
}

// SelectFromBufferedEvents reads one or more rows from buffered_events table
func (pdb *db) SelectFromBufferedEvents(ctx context.Context, filter *sqlplugin.BufferedEventsFilter) ([]sqlplugin.BufferedEventsRow, error) {
	var rows []sqlplugin.BufferedEventsRow
	err := pdb.conn.SelectContext(ctx, &rows, getBufferedEventsQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
	for i := 0; i < len(rows); i++ {
		rows[i].DomainID = filter.DomainID
		rows[i].WorkflowID = filter.WorkflowID
		rows[i].RunID = filter.RunID
		rows[i].ShardID = filter.ShardID
	}
	return rows, err
}

// DeleteFromBufferedEvents deletes one or more rows from buffered_events table
func (pdb *db) DeleteFromBufferedEvents(ctx context.Context, filter *sqlplugin.BufferedEventsFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteBufferedEventsQuery, filter.ShardID, filter.DomainID, filter.WorkflowID, filter.RunID)
}

// InsertIntoReplicationTasks inserts one or more rows into replication_tasks table
func (pdb *db) InsertIntoReplicationTasks(ctx context.Context, rows []sqlplugin.ReplicationTasksRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, createReplicationTasksQuery, rows)
}

// SelectFromReplicationTasks reads one or more rows from replication_tasks table
func (pdb *db) SelectFromReplicationTasks(ctx context.Context, filter *sqlplugin.ReplicationTasksFilter) ([]sqlplugin.ReplicationTasksRow, error) {
	var rows []sqlplugin.ReplicationTasksRow
	err := pdb.conn.SelectContext(ctx, &rows, getReplicationTasksQuery, filter.ShardID, filter.MinTaskID, filter.MaxTaskID, filter.PageSize)
	return rows, err
}

// DeleteFromReplicationTasks deletes one rows from replication_tasks table
func (pdb *db) DeleteFromReplicationTasks(ctx context.Context, filter *sqlplugin.ReplicationTasksFilter) (sql.Result, error) {
	return pdb.conn.ExecContext(ctx, deleteReplicationTaskQuery, filter.ShardID, filter.TaskID)
}

// RangeDeleteFromReplicationTasks deletes multi rows from replication_tasks table
func (pdb *db) RangeDeleteFromReplicationTasks(ctx context.Context, filter *sqlplugin.ReplicationTasksFilter) (sql.Result, error) {
	if filter.PageSize > 0 {
		return pdb.conn.ExecContext(ctx, rangeDeleteReplicationTaskByBatchQuery, filter.ShardID, filter.InclusiveEndTaskID, filter.PageSize)
	}
	return pdb.conn.ExecContext(ctx, rangeDeleteReplicationTaskQuery, filter.ShardID, filter.InclusiveEndTaskID)
}

// InsertIntoReplicationTasksDLQ inserts one or more rows into replication_tasks_dlq table
func (pdb *db) InsertIntoReplicationTasksDLQ(ctx context.Context, row *sqlplugin.ReplicationTaskDLQRow) (sql.Result, error) {
	return pdb.conn.NamedExecContext(ctx, insertReplicationTaskDLQQuery, row)
}

// SelectFromReplicationTasksDLQ reads one or more rows from replication_tasks_dlq table
func (pdb *db) SelectFromReplicationTasksDLQ(ctx context.Context, filter *sqlplugin.ReplicationTasksDLQFilter) ([]sqlplugin.ReplicationTasksRow, error) {
	var rows []sqlplugin.ReplicationTasksRow
	err := pdb.conn.SelectContext(
		ctx,
		&rows, getReplicationTasksDLQQuery,
		filter.SourceClusterName,
		filter.ShardID,
		filter.MinTaskID,
		filter.MaxTaskID,
		filter.PageSize)
	return rows, err
}

// SelectFromReplicationDLQ reads one row from replication_tasks_dlq table
func (pdb *db) SelectFromReplicationDLQ(ctx context.Context, filter *sqlplugin.ReplicationTaskDLQFilter) (int64, error) {
	var size []int64
	if err := pdb.conn.SelectContext(
		ctx,
		&size, getReplicationTaskDLQQuery,
		filter.SourceClusterName,
		filter.ShardID,
	); err != nil {
		return 0, err
	}
	return size[0], nil
}

// DeleteMessageFromReplicationTasksDLQ deletes one row from replication_tasks_dlq table
func (pdb *db) DeleteMessageFromReplicationTasksDLQ(
	ctx context.Context,
	filter *sqlplugin.ReplicationTasksDLQFilter,
) (sql.Result, error) {

	return pdb.conn.ExecContext(
		ctx,
		deleteReplicationTaskFromDLQQuery,
		filter.SourceClusterName,
		filter.ShardID,
		filter.TaskID,
	)
}

// DeleteMessageFromReplicationTasksDLQ deletes one or more rows from replication_tasks_dlq table
func (pdb *db) RangeDeleteMessageFromReplicationTasksDLQ(
	ctx context.Context,
	filter *sqlplugin.ReplicationTasksDLQFilter,
) (sql.Result, error) {
	if filter.PageSize > 0 {
		return pdb.conn.ExecContext(
			ctx,
			rangeDeleteReplicationTaskFromDLQByBatchQuery,
			filter.SourceClusterName,
			filter.ShardID,
			filter.TaskID,
			filter.InclusiveEndTaskID,
			filter.PageSize,
		)
	}

	return pdb.conn.ExecContext(
		ctx,
		rangeDeleteReplicationTaskFromDLQQuery,
		filter.SourceClusterName,
		filter.ShardID,
		filter.TaskID,
		filter.InclusiveEndTaskID,
	)
}

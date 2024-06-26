// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Conversation is an object representing the database table.
type Conversation struct {
	UUID      string    `boil:"uuid" json:"uuid" toml:"uuid" yaml:"uuid"`
	CreatedAt null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *conversationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L conversationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ConversationColumns = struct {
	UUID      string
	CreatedAt string
	UpdatedAt string
}{
	UUID:      "uuid",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ConversationTableColumns = struct {
	UUID      string
	CreatedAt string
	UpdatedAt string
}{
	UUID:      "conversations.uuid",
	CreatedAt: "conversations.created_at",
	UpdatedAt: "conversations.updated_at",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod   { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod  { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) ILIKE(x string) qm.QueryMod  { return qm.Where(w.field+" ILIKE ?", x) }
func (w whereHelperstring) NILIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT ILIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var ConversationWhere = struct {
	UUID      whereHelperstring
	CreatedAt whereHelpernull_Time
	UpdatedAt whereHelpernull_Time
}{
	UUID:      whereHelperstring{field: "\"conversations\".\"uuid\""},
	CreatedAt: whereHelpernull_Time{field: "\"conversations\".\"created_at\""},
	UpdatedAt: whereHelpernull_Time{field: "\"conversations\".\"updated_at\""},
}

// ConversationRels is where relationship names are stored.
var ConversationRels = struct {
	Messages string
}{
	Messages: "Messages",
}

// conversationR is where relationships are stored.
type conversationR struct {
	Messages MessageSlice `boil:"Messages" json:"Messages" toml:"Messages" yaml:"Messages"`
}

// NewStruct creates a new relationship struct
func (*conversationR) NewStruct() *conversationR {
	return &conversationR{}
}

func (r *conversationR) GetMessages() MessageSlice {
	if r == nil {
		return nil
	}
	return r.Messages
}

// conversationL is where Load methods for each relationship are stored.
type conversationL struct{}

var (
	conversationAllColumns            = []string{"uuid", "created_at", "updated_at"}
	conversationColumnsWithoutDefault = []string{"uuid"}
	conversationColumnsWithDefault    = []string{"created_at", "updated_at"}
	conversationPrimaryKeyColumns     = []string{"uuid"}
	conversationGeneratedColumns      = []string{}
)

type (
	// ConversationSlice is an alias for a slice of pointers to Conversation.
	// This should almost always be used instead of []Conversation.
	ConversationSlice []*Conversation
	// ConversationHook is the signature for custom Conversation hook methods
	ConversationHook func(context.Context, boil.ContextExecutor, *Conversation) error

	conversationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	conversationType                 = reflect.TypeOf(&Conversation{})
	conversationMapping              = queries.MakeStructMapping(conversationType)
	conversationPrimaryKeyMapping, _ = queries.BindMapping(conversationType, conversationMapping, conversationPrimaryKeyColumns)
	conversationInsertCacheMut       sync.RWMutex
	conversationInsertCache          = make(map[string]insertCache)
	conversationUpdateCacheMut       sync.RWMutex
	conversationUpdateCache          = make(map[string]updateCache)
	conversationUpsertCacheMut       sync.RWMutex
	conversationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var conversationAfterSelectMu sync.Mutex
var conversationAfterSelectHooks []ConversationHook

var conversationBeforeInsertMu sync.Mutex
var conversationBeforeInsertHooks []ConversationHook
var conversationAfterInsertMu sync.Mutex
var conversationAfterInsertHooks []ConversationHook

var conversationBeforeUpdateMu sync.Mutex
var conversationBeforeUpdateHooks []ConversationHook
var conversationAfterUpdateMu sync.Mutex
var conversationAfterUpdateHooks []ConversationHook

var conversationBeforeDeleteMu sync.Mutex
var conversationBeforeDeleteHooks []ConversationHook
var conversationAfterDeleteMu sync.Mutex
var conversationAfterDeleteHooks []ConversationHook

var conversationBeforeUpsertMu sync.Mutex
var conversationBeforeUpsertHooks []ConversationHook
var conversationAfterUpsertMu sync.Mutex
var conversationAfterUpsertHooks []ConversationHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Conversation) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Conversation) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Conversation) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Conversation) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Conversation) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Conversation) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Conversation) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Conversation) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Conversation) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddConversationHook registers your hook function for all future operations.
func AddConversationHook(hookPoint boil.HookPoint, conversationHook ConversationHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		conversationAfterSelectMu.Lock()
		conversationAfterSelectHooks = append(conversationAfterSelectHooks, conversationHook)
		conversationAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		conversationBeforeInsertMu.Lock()
		conversationBeforeInsertHooks = append(conversationBeforeInsertHooks, conversationHook)
		conversationBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		conversationAfterInsertMu.Lock()
		conversationAfterInsertHooks = append(conversationAfterInsertHooks, conversationHook)
		conversationAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		conversationBeforeUpdateMu.Lock()
		conversationBeforeUpdateHooks = append(conversationBeforeUpdateHooks, conversationHook)
		conversationBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		conversationAfterUpdateMu.Lock()
		conversationAfterUpdateHooks = append(conversationAfterUpdateHooks, conversationHook)
		conversationAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		conversationBeforeDeleteMu.Lock()
		conversationBeforeDeleteHooks = append(conversationBeforeDeleteHooks, conversationHook)
		conversationBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		conversationAfterDeleteMu.Lock()
		conversationAfterDeleteHooks = append(conversationAfterDeleteHooks, conversationHook)
		conversationAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		conversationBeforeUpsertMu.Lock()
		conversationBeforeUpsertHooks = append(conversationBeforeUpsertHooks, conversationHook)
		conversationBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		conversationAfterUpsertMu.Lock()
		conversationAfterUpsertHooks = append(conversationAfterUpsertHooks, conversationHook)
		conversationAfterUpsertMu.Unlock()
	}
}

// One returns a single conversation record from the query.
func (q conversationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Conversation, error) {
	o := &Conversation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for conversations")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Conversation records from the query.
func (q conversationQuery) All(ctx context.Context, exec boil.ContextExecutor) (ConversationSlice, error) {
	var o []*Conversation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Conversation slice")
	}

	if len(conversationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Conversation records in the query.
func (q conversationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count conversations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q conversationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if conversations exists")
	}

	return count > 0, nil
}

// Messages retrieves all the message's Messages with an executor.
func (o *Conversation) Messages(mods ...qm.QueryMod) messageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"messages\".\"conversation_id\"=?", o.UUID),
	)

	return Messages(queryMods...)
}

// LoadMessages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (conversationL) LoadMessages(ctx context.Context, e boil.ContextExecutor, singular bool, maybeConversation interface{}, mods queries.Applicator) error {
	var slice []*Conversation
	var object *Conversation

	if singular {
		var ok bool
		object, ok = maybeConversation.(*Conversation)
		if !ok {
			object = new(Conversation)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeConversation)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeConversation))
			}
		}
	} else {
		s, ok := maybeConversation.(*[]*Conversation)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeConversation)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeConversation))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &conversationR{}
		}
		args[object.UUID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &conversationR{}
			}
			args[obj.UUID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`messages`),
		qm.WhereIn(`messages.conversation_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load messages")
	}

	var resultSlice []*Message
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice messages")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on messages")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for messages")
	}

	if len(messageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Messages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &messageR{}
			}
			foreign.R.Conversation = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.UUID, foreign.ConversationID) {
				local.R.Messages = append(local.R.Messages, foreign)
				if foreign.R == nil {
					foreign.R = &messageR{}
				}
				foreign.R.Conversation = local
				break
			}
		}
	}

	return nil
}

// AddMessages adds the given related objects to the existing relationships
// of the conversation, optionally inserting them as new records.
// Appends related to o.R.Messages.
// Sets related.R.Conversation appropriately.
func (o *Conversation) AddMessages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Message) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.ConversationID, o.UUID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"messages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"conversation_id"}),
				strmangle.WhereClause("\"", "\"", 2, messagePrimaryKeyColumns),
			)
			values := []interface{}{o.UUID, rel.UUID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.ConversationID, o.UUID)
		}
	}

	if o.R == nil {
		o.R = &conversationR{
			Messages: related,
		}
	} else {
		o.R.Messages = append(o.R.Messages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &messageR{
				Conversation: o,
			}
		} else {
			rel.R.Conversation = o
		}
	}
	return nil
}

// SetMessages removes all previously related items of the
// conversation replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Conversation's Messages accordingly.
// Replaces o.R.Messages with related.
// Sets related.R.Conversation's Messages accordingly.
func (o *Conversation) SetMessages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Message) error {
	query := "update \"messages\" set \"conversation_id\" = null where \"conversation_id\" = $1"
	values := []interface{}{o.UUID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Messages {
			queries.SetScanner(&rel.ConversationID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Conversation = nil
		}
		o.R.Messages = nil
	}

	return o.AddMessages(ctx, exec, insert, related...)
}

// RemoveMessages relationships from objects passed in.
// Removes related items from R.Messages (uses pointer comparison, removal does not keep order)
// Sets related.R.Conversation.
func (o *Conversation) RemoveMessages(ctx context.Context, exec boil.ContextExecutor, related ...*Message) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.ConversationID, nil)
		if rel.R != nil {
			rel.R.Conversation = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("conversation_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Messages {
			if rel != ri {
				continue
			}

			ln := len(o.R.Messages)
			if ln > 1 && i < ln-1 {
				o.R.Messages[i] = o.R.Messages[ln-1]
			}
			o.R.Messages = o.R.Messages[:ln-1]
			break
		}
	}

	return nil
}

// Conversations retrieves all the records using an executor.
func Conversations(mods ...qm.QueryMod) conversationQuery {
	mods = append(mods, qm.From("\"conversations\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"conversations\".*"})
	}

	return conversationQuery{q}
}

// FindConversation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindConversation(ctx context.Context, exec boil.ContextExecutor, uUID string, selectCols ...string) (*Conversation, error) {
	conversationObj := &Conversation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"conversations\" where \"uuid\"=$1", sel,
	)

	q := queries.Raw(query, uUID)

	err := q.Bind(ctx, exec, conversationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from conversations")
	}

	if err = conversationObj.doAfterSelectHooks(ctx, exec); err != nil {
		return conversationObj, err
	}

	return conversationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Conversation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no conversations provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(conversationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	conversationInsertCacheMut.RLock()
	cache, cached := conversationInsertCache[key]
	conversationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			conversationAllColumns,
			conversationColumnsWithDefault,
			conversationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(conversationType, conversationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"conversations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"conversations\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into conversations")
	}

	if !cached {
		conversationInsertCacheMut.Lock()
		conversationInsertCache[key] = cache
		conversationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Conversation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Conversation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	conversationUpdateCacheMut.RLock()
	cache, cached := conversationUpdateCache[key]
	conversationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			conversationAllColumns,
			conversationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update conversations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"conversations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, conversationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, append(wl, conversationPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update conversations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for conversations")
	}

	if !cached {
		conversationUpdateCacheMut.Lock()
		conversationUpdateCache[key] = cache
		conversationUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q conversationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for conversations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for conversations")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ConversationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"conversations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, conversationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in conversation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all conversation")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Conversation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no conversations provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(conversationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	conversationUpsertCacheMut.RLock()
	cache, cached := conversationUpsertCache[key]
	conversationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			conversationAllColumns,
			conversationColumnsWithDefault,
			conversationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			conversationAllColumns,
			conversationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert conversations, could not build update column list")
		}

		ret := strmangle.SetComplement(conversationAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(conversationPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert conversations, could not build conflict column list")
			}

			conflict = make([]string, len(conversationPrimaryKeyColumns))
			copy(conflict, conversationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"conversations\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(conversationType, conversationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert conversations")
	}

	if !cached {
		conversationUpsertCacheMut.Lock()
		conversationUpsertCache[key] = cache
		conversationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Conversation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Conversation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Conversation provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), conversationPrimaryKeyMapping)
	sql := "DELETE FROM \"conversations\" WHERE \"uuid\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from conversations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for conversations")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q conversationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no conversationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from conversations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for conversations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ConversationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(conversationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"conversations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, conversationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from conversation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for conversations")
	}

	if len(conversationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Conversation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindConversation(ctx, exec, o.UUID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ConversationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ConversationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"conversations\".* FROM \"conversations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, conversationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ConversationSlice")
	}

	*o = slice

	return nil
}

// ConversationExists checks if the Conversation row exists.
func ConversationExists(ctx context.Context, exec boil.ContextExecutor, uUID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"conversations\" where \"uuid\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, uUID)
	}
	row := exec.QueryRowContext(ctx, sql, uUID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if conversations exists")
	}

	return exists, nil
}

// Exists checks if the Conversation row exists.
func (o *Conversation) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ConversationExists(ctx, exec, o.UUID)
}

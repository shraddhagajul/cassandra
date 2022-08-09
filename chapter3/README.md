# cassandra
cassandra-definitive-guide

# Introduction
1. Cassandra defines a column family to be a logical division that associates similar data.
2. Both row keys and column names can be strings, like relational column names, but they can also be long integers, UUIDs, or any kind of byte array.
3. Instead of storing null for those values we don’t know, which would waste space, we just won’t store that column at all for that row. So now we have a sparse, multidimensional array structure
4. Columns in Cassandra actually have a third aspect: the timestamp, which records the last time the column was updated.
   You cannot query by the timestamp. Rows do not have timestamps. Only each individual column has a timestamp.
5. A super column family can be thought of as a map of maps. Create a group of related columns.
   Where a row in a column family holds a collection of name/value pairs, the super column family holds subcolumns, where subcolumns are named groups of columns.
   a row in a super column family still contains columns, each of which then contains subcolumns.

# Clusters
1. The outermost structure in Cassandra is the cluster, sometimes called the ring, because Cassandra assigns data to nodes in the cluster by arranging them in a ring.
2. A node holds a replica for different ranges of data
3. If the first node goes down, a replica can respond to queries. The peer-to-peer protocol allows the data to replicate across nodes in a manner transparent to the user.
4. Replication factor is the number of machines in your cluster that will receive copies of the same data.

# Keyspaces
1. A keyspace is the outermost container for a list of one or more column families.
2. Like a relational database, a keyspace has a name and a set of attributes that define keyspace-wide behavior.
3. Depending on your security constraints and partitioner, it’s fine to run multiple key- spaces on the same cluster. 
   For example, if your application is called Twitter, you would probably have a cluster called Twitter-Cluster and a keyspace called Twitter.
4. The only time you would want to split your application into multiple keyspaces is if you wanted a different replication factor or replica placement strategy for some of the column families.
5. The basic attributes that you can set per keyspace are : 
   1. The replication factor 
      1. refers to the number of nodes that will act as copies (replicas) of each row of data.
      2. essentially allows you to decide how much you want to pay in performance to gain more consistency. 
      3. That is, your consistency level for reading and writing data is based on the replication factor.
   2. Replica placement strategy
      1. The replica placement refers to how the replicas will be placed in the ring.
      2. There are different strategies that ship with Cassandra for determining which nodes will get copies of which keys.
      3. eg : SimpleStrategy, OldNetworkTopologyStrategy, NetworkTopologyStrategy
   3. Column families
      1. A column family is roughly analagous to a table in the relational model, and is a container for a collection of rows, each of which is itself an ordered collection of columns.
      2. Column families represent the structure of your data.
      3. Each keyspace has at least one and often many column families.

# Column Families
   1. RDMS vs Cassandra
      1. Cassandra is considered schema-free because although the column families are defined, the columns are not. You can freely add any column to any column family at any time, depending on your needs.
      2. a column family has two attributes: a name and a comparator. The comparator value indicates how columns will be sorted when they are returned to you in a query
      3. it is rare to hear of recommendations about data modeling based on how the RDBMS might store tables on disk. 
         That’s another reason to keep in mind that a column family is not a table. Because column families are each stored in separate files on disk, it’s important to keep related columns defined together in the same column family.
      4. a table can hold columns, or it can be defined as a super column family. The benefit of using a super column family is to allow for nesting.
      5. in Cassandra, you specify values for one or more columns. That collection of values together with a unique identifier is called a row. That row has a unique key, called the row key, which acts like the primary key unique identifier for that row. 
      6. think of rows as containers for columns. some people refer to Cassandra column families as similar to a four-dimensional hash:
         [Keyspace][ColumnFamily][Key][Column]

# FYI
   It’s an inherent part of Cassandra’s replica design that all data for a single row must fit on a single machine in the cluster. 
   The reason for this limitation is that rows have an associated row key, which is used to determine the nodes that will act as replicas for that row. 
   Further, the value of a single column cannot exceed 2GB. 

# Column Families options
   1. keys_cached
      1. The number of locations to keep cached per SSTable. 
      2. This doesn’t refer to column name/values at all, but to the number of keys, as locations of rows per column family, to keep in memory in least-recently-used order.
   2. rows_cached
      The number of rows whose entire contents will be cached in memory.
   3. comment
      column family definitions.
   4. read_repair_chance
   5. preload_row_cache
      Specifies whether you want to prepopulate the row cache on server startup.

# Column
    A column is a triplet of a name, a value, and a clock
    1. in Cassandra, you don’t define the columns up front; you just define the column families you want in the keyspace, and then you can start writing data without defining the columns anywhere. 
    2. That’s because in Cassandra, all of a column’s names are supplied by the client. 
    3. This adds considerable flexibility to how your application works with data, and can allow it to evolve organically over time.
    4. we think of a relational table as holding the same set of columns for every row. But in Cassandra, a column family holds many rows, each of which may hold the same, or different, sets of columns.
    5. On the server side, columns are immutable in order to prevent multithreading issues.
    6. rows for the same column family are stored together on disk.

# wide and skinny model 
    1. A wide row means a row that has lots and lots of columns.
    2. smaller number of columns and use many different rows—that’s the skinny model.
    3. Skinny rows are slightly more like traditional RDBMS rows, but all columns are optional
    4. Another difference between wide and skinny rows is that only wide rows will typically be concerned about sorting order of column names.

# Column sorting
    1. In Cassandra, you specify how column names will be compared for sort order when results are returned to the client.
    2. Columns are sorted by the “Compare With” type defined on their enclosing column family
    3. It is not possible in Cassandra to sort by value
    4. This may seem like an odd limitation, but Cassandra has to sort by column name in order to allow fetching individual columns from very large rows without pulling the entire row into memory. 
        Performance is an important selling point of Cassandra, and sorting at read time would harm performance.

# Super Column
    1. The basic structure of a super column is its name, which is a byte array (just as with a regular column), and the columns it stores.
    2. Its columns are held as a map whose keys are the column names and whose values are the columns.
    3. for super columns, it becomes more like a five-dimensional hash:
        [Keyspace][ColumnFamily][Key][SuperColumn][SubColumn]
    4. To use a super column, you define your column family as type Super. 
    5. Then, you still have row keys as you do in a regular column family, but you also reference the super column, which is simply a name that points to a list or map of regular columns (some- times called the subcolumns).
# Materialized View

# Valueless Column

# how Cassandra is different from RDBMS
    1. No Query Language
    2. No Referential Integrity
    3. Secondary Indexes

# Materialized View
    1. “materialized” means storing a full copy of the original data so that everything you need to answer a query is right there, without forcing you to look up the original data.
    2. If you are performing a second query because you’re only storing column names that you use, like foreign keys in the second column family, that’s a secondary index.

# Valueless Column

# Aggregate Key

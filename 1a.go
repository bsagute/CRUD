func TestGetNewLayoutData(t *testing.T) {
    // Create a new sqlmock database connection and a mock object
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    // Open a new GORM DB connection using the mock database
    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    assert.NoError(t, err)

    // Positive Case: First Page (Parent Entity ID is empty)
    t.Run("FirstPage", func(t *testing.T) {
        parentEntityID := ""

        // Mock the entity type ID query
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs("00000000-0000-0000-0000-000000000000").
            WillReturnRows(sqlmock.NewRows([]string{"entity_id"}).AddRow(uuid.New()))

        // Mock the entity info query
        mock.ExpectQuery(`SELECT * FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs("00000000-0000-0000-0000-000000000000").
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "has_child"}).AddRow(uuid.New(), true))

        // Mock the metric layout info query
        mock.ExpectQuery(`SELECT .* FROM metric m JOIN entity_info e ON m.entity_id = e.entity_id WHERE m.parent_entity_id = \$1 GROUP BY m.entity_id, m.metric_id, sort ORDER BY sort ASC`).
            WithArgs("00000000-0000-0000-0000-000000000000").
            WillReturnRows(sqlmock.NewRows([]string{"metric_id", "text", "parent_entity_id", "entity_name", "entity_type", "has_child", "deepLink_id", "metric_name", "metric_type", "metric_unit", "min_threshold", "max_threshold"}).
                AddRow("metric_id_1", "metric_text_1", "00000000-0000-0000-0000-000000000000", "entity_name_1", "entity_type_1", true, "deepLink_id_1", "metric_name_1", `[{"graphType": {"colSpan": "6", "sort": "1", "pageSort": "1", "containerClass": "class_1", "isCollapsible": true, "containerName": "name_1", "pageClass": "page_class_1", "pageTitle": "page_title_1", "metricType": "metric_type_1"}}]`, "metric_unit_1", 10, 20))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.NoError(t, err)
        assert.NotNil(t, result)

        // Assert the content of the result
        assert.Equal(t, 1, len(result.Layout))
        assert.Equal(t, "metric_id_1", result.Layout[0].MetricId)
        assert.Equal(t, "metric_text_1", result.Layout[0].Text)
        assert.Equal(t, "00000000-0000-0000-0000-000000000000", result.Layout[0].ParentEntityId)
        assert.Equal(t, "entity_name_1", result.Layout[0].EntityName)
        assert.Equal(t, "entity_type_1", result.Layout[0].EntityType)
        assert.Equal(t, true, result.Layout[0].HasChild)
        assert.Equal(t, "deepLink_id_1", result.Layout[0].DeepLinkId)
        assert.Equal(t, "metric_name_1", result.Layout[0].MetricName)
        assert.Equal(t, "metric_type_1", result.Layout[0].MetricType)
        assert.Equal(t, "metric_unit_1", result.Layout[0].MetricUnit)
        assert.Equal(t, 10, result.Layout[0].MinThreshold)
        assert.Equal(t, 20, result.Layout[0].MaxThreshold)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    // Negative Case: Error during Entity Type ID Query
    t.Run("ErrorEntityTypeIDQuery", func(t *testing.T) {
        parentEntityID := "00000000-0000-0000-0000-000000000000"

        // Mock the entity type ID query to return an error
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs(parentEntityID).
            WillReturnError(errors.New("query error"))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.Error(t, err)
        assert.Empty(t, result.Layout)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    // Negative Case: Error during Entity Info Query
    t.Run("ErrorEntityInfoQuery", func(t *testing.T) {
        parentEntityID := "00000000-0000-0000-0000-000000000000"

        // Mock the entity type ID query
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id"}).AddRow(uuid.New()))

        // Mock the entity info query to return an error
        mock.ExpectQuery(`SELECT * FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(parentEntityID).
            WillReturnError(errors.New("query error"))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.Error(t, err)
        assert.Empty(t, result.Layout)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    // Negative Case: Error during Metric Layout Info Query
    t.Run("ErrorMetricLayoutInfoQuery", func(t *testing.T) {
        parentEntityID := "00000000-0000-0000-0000-000000000000"

        // Mock the entity type ID query
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id"}).AddRow(uuid.New()))

        // Mock the entity info query
        mock.ExpectQuery(`SELECT * FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "has_child"}).AddRow(uuid.New(), true))

        // Mock the metric layout info query to return an error
        mock.ExpectQuery(`SELECT .* FROM metric m JOIN entity_info e ON m.entity_id = e.entity_id WHERE m.parent_entity_id = \$1 GROUP BY m.entity_id, m.metric_id, sort ORDER BY sort ASC`).
            WithArgs(parentEntityID).
            WillReturnError(errors.New("query error"))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.Error(t, err)
        assert.Empty(t, result.Layout)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    // Edge Case: Empty Result Set
    t.Run("EmptyResultSet", func(t *testing.T) {
        parentEntityID := "00000000-0000-0000-0000-000000000000"

        // Mock the entity type ID query
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id"}).AddRow(uuid.New()))

        // Mock the entity info query
        mock.ExpectQuery(`SELECT * FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "has_child"}).AddRow(uuid.New(), true))

        // Mock the metric layout info query to return an empty result set
        mock.ExpectQuery(`SELECT .* FROM metric m JOIN entity_info e ON m.entity_id = e.entity_id WHERE m.parent_entity_id = \$1 GROUP BY m.entity_id, m.metric_id, sort ORDER BY sort ASC`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"metric_id", "text", "parent_entity_id", "entity_name", "entity_type", "has_child", "deepLink_id", "metric_name", "metric_type", "metric_unit", "min_threshold", "max_threshold"}))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.NoError(t, err)
        assert.NotNil(t, result



assert.NotNil(t, result)
        assert.Empty(t, result.Layout)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    // Positive Case: Normal flow with non-empty parent_entity_id
    t.Run("NormalFlow", func(t *testing.T) {
        parentEntityID := "123e4567-e89b-12d3-a456-426614174000"

        // Mock the entity type ID query
        mock.ExpectQuery(`SELECT DISTINCT "entity_info"."entity_id" FROM "entity_info" WHERE parent_entity_id = \$1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id"}).AddRow(uuid.New()))

        // Mock the entity info query
        mock.ExpectQuery(`SELECT * FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "has_child"}).AddRow(uuid.New(), true))

        // Mock the metric layout info query
        mock.ExpectQuery(`SELECT .* FROM metric m JOIN entity_info e ON m.entity_id = e.entity_id WHERE m.parent_entity_id = \$1 GROUP BY m.entity_id, m.metric_id, sort ORDER BY sort ASC`).
            WithArgs(parentEntityID).
            WillReturnRows(sqlmock.NewRows([]string{"metric_id", "text", "parent_entity_id", "entity_name", "entity_type", "has_child", "deepLink_id", "metric_name", "metric_type", "metric_unit", "min_threshold", "max_threshold"}).
                AddRow("metric_id_1", "metric_text_1", parentEntityID, "entity_name_1", "entity_type_1", true, "deepLink_id_1", "metric_name_1", `[{"graphType": {"colSpan": "6", "sort": "1", "pageSort": "1", "containerClass": "class_1", "isCollapsible": true, "containerName": "name_1", "pageClass": "page_class_1", "pageTitle": "page_title_1", "metricType": "metric_type_1"}}]`, "metric_unit_1", 10, 20))

        // Call the function
        result, err := db.GetNewLayoutData(parentEntityID, gdb)
        assert.NoError(t, err)
        assert.NotNil(t, result)

        // Assert the content of the result
        assert.Equal(t, 1, len(result.Layout))
        assert.Equal(t, "metric_id_1", result.Layout[0].MetricId)
        assert.Equal(t, "metric_text_1", result.Layout[0].Text)
        assert.Equal(t, parentEntityID, result.Layout[0].ParentEntityId)
        assert.Equal(t, "entity_name_1", result.Layout[0].EntityName)
        assert.Equal(t, "entity_type_1", result.Layout[0].EntityType)
        assert.Equal(t, true, result.Layout[0].HasChild)
        assert.Equal(t, "deepLink_id_1", result.Layout[0].DeepLinkId)
        assert.Equal(t, "metric_name_1", result.Layout[0].MetricName)
        assert.Equal(t, "metric_type_1", result.Layout[0].MetricType)
        assert.Equal(t, "metric_unit_1", result.Layout[0].MetricUnit)
        assert.Equal(t, 10, result.Layout[0].MinThreshold)
        assert.Equal(t, 20, result.Layout[0].MaxThreshold)

        // Verify that all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}


	  

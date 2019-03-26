package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TemplateType for messages
type TemplateType struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	Name         string                 `bson:"name"`
	Translations map[string]Translation `bson:"translations"` // guaranteed to include english translation
}

func (t TemplateType) addTranslation(language string, translation Translation) (err error) {
	if _, exists := t.Translations[language]; exists {
		return errors.New("translation already exists")
	}
	t.Translations[language] = translation

	return
} //

func (t TemplateType) updateTranslation(language string, translation Translation) (err error) {
	if _, exists := t.Translations[language]; !exists {
		return errors.New("translation does not exist")
	}
	t.Translations[language] = translation

	return
} //

func (t TemplateType) getTranslationOrDefault(language string) (translation Translation) {
	var exists bool
	translation, exists = t.Translations[language]
	if !exists {
		translation, _ = t.Translations["EN"] //default
	}

	return //
}

// Translation of TemplateType
type Translation struct {
	HTML    string   `bson:"html"`    // html version of template
	Text    string   `bson:"text"`    // text version of template
	Expects []string `bson:"expects"` // parameters needed from e.g. the users profile
}

/*
######## DB Methods ########
*/

func getCollection() *mongo.Collection {
	return dbClient.Database("template-prototype").Collection("templates")
} //

func getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func createTemplate(name string, translation Translation) (err error) {
	_, err = getTemplateTypeByName(name)
	if err == nil {
		err = errors.New("template exists")
		return
	}

	t := TemplateType{
		Name:         name,
		Translations: map[string]Translation{},
	}
	t.Translations["EN"] = translation

	ctx, cancel := getContext()
	defer cancel()

	_, err = getCollection().InsertOne(ctx, t)
	if err != nil {
		return
	}

	return
} //

func createTranslation(name string, language string, translation Translation) (err error) {
	t, err := getTemplateTypeByName(name)
	if err != nil {
		return
	}

	err = t.addTranslation(language, translation)
	if err != nil {
		return
	}

	err = updateTemplateType(t)

	return
} //

func updateTranslation(name string, language string, translation Translation) (err error) {
	t, err := getTemplateTypeByName(name)
	if err != nil {
		return
	}

	err = t.updateTranslation(language, translation)

	return
}

func getTranslationOrDefault(name string, language string) (translation Translation, err error) {
	t, err := getTemplateTypeByName(name)
	if err != nil {
		return
	}

	translation = t.getTranslationOrDefault(language)

	return
}

func getTemplateTypeByName(name string) (t TemplateType, err error) {
	ctx, cancel := getContext()
	defer cancel()

	t = TemplateType{}
	filter := bson.M{"name": name}
	err = getCollection().FindOne(ctx, filter).Decode(&t)

	return
} //

func updateTemplateType(t TemplateType) (err error) {
	ctx, cancel := getContext()
	defer cancel()

	tOld := TemplateType{}
	filter := bson.M{"_id": t.ID}

	err = getCollection().FindOneAndReplace(ctx, filter, t).Decode(&tOld)

	return
} //
